package providers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/erdemkosk/gitness/internal/util"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubProvider struct {
	client *githubv4.Client
}

func NewGitHubProvider(token string) *GitHubProvider {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)
	return &GitHubProvider{
		client: githubv4.NewClient(httpClient),
	}
}

func (g *GitHubProvider) FetchCommits(owner, repo string, duration string) (map[string]CommitInfo, error) {
	const (
		maxWorkers = 10
		batchSize  = 100
	)

	type commitBatch struct {
		commits []struct {
			Author struct {
				Name  string
				Email string
			}
			CommittedDate githubv4.DateTime
		}
		err error
	}

	var since *githubv4.GitTimestamp
	if duration != "" {
		dur, err := util.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %w", err)
		}
		since = &githubv4.GitTimestamp{Time: dur.ToTime()}
	}

	// Initial query to get cursor
	var initialQuery struct {
		Repository struct {
			DefaultBranchRef struct {
				Target struct {
					Commit struct {
						History struct {
							TotalCount int
							PageInfo   struct {
								HasNextPage bool
								EndCursor   string
							}
						} `graphql:"history(first: $limit, since: $since, after: $after)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"since": since,
		"limit": githubv4.Int(1),
		"after": (*githubv4.String)(nil),
	}

	if err := g.client.Query(context.Background(), &initialQuery, variables); err != nil {
		return nil, fmt.Errorf("initial GitHub query failed: %v", err)
	}

	totalCommits := initialQuery.Repository.DefaultBranchRef.Target.Commit.History.TotalCount
	if totalCommits == 0 {
		return make(map[string]CommitInfo), nil
	}

	// Prepare workers for parallel execution
	results := make(chan commitBatch, totalCommits/batchSize+1)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)
	var mu sync.Mutex
	cursors := []string{""}

	if initialQuery.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage {
		cursors = append(cursors, initialQuery.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
	}

	// Collect cursors
	for len(cursors) < totalCommits/batchSize+1 {
		var query struct {
			Repository struct {
				DefaultBranchRef struct {
					Target struct {
						Commit struct {
							History struct {
								PageInfo struct {
									HasNextPage bool
									EndCursor   string
								}
							} `graphql:"history(first: $limit, after: $after, since: $since)"`
						} `graphql:"... on Commit"`
					}
				}
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		variables["after"] = githubv4.String(cursors[len(cursors)-1])
		variables["limit"] = githubv4.Int(batchSize)

		if err := g.client.Query(context.Background(), &query, variables); err != nil {
			return nil, fmt.Errorf("cursor collection failed: %v", err)
		}

		if !query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage {
			break
		}

		cursors = append(cursors, query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
	}

	// Parallel data fetching
	for i, cursor := range cursors {
		wg.Add(1)
		semaphore <- struct{}{} // Rate limiting

		go func(after string, batchNum int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			var query struct {
				Repository struct {
					DefaultBranchRef struct {
						Target struct {
							Commit struct {
								History struct {
									Nodes []struct {
										Author struct {
											Name  string
											Email string
										}
										CommittedDate githubv4.DateTime
									}
								} `graphql:"history(first: $limit, after: $after, since: $since)"`
							} `graphql:"... on Commit"`
						}
					}
				} `graphql:"repository(owner: $owner, name: $repo)"`
			}

			afterCursor := (*githubv4.String)(nil)
			if after != "" {
				temp := githubv4.String(after)
				afterCursor = &temp
			}

			variables := map[string]interface{}{
				"owner": githubv4.String(owner),
				"repo":  githubv4.String(repo),
				"limit": githubv4.Int(batchSize),
				"after": afterCursor,
				"since": since,
			}

			err := g.client.Query(context.Background(), &query, variables)

			mu.Lock()
			results <- commitBatch{
				commits: query.Repository.DefaultBranchRef.Target.Commit.History.Nodes,
				err:     err,
			}
			mu.Unlock()

			// Rate limiting
			time.Sleep(time.Second / 5)
		}(cursor, i)
	}

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Merge results
	authorStats := make(map[string]CommitInfo)
	for batch := range results {
		if batch.err != nil {
			return nil, fmt.Errorf("batch query failed: %v", batch.err)
		}

		for _, commit := range batch.commits {
			author := commit.Author.Name
			if author == "" {
				author = commit.Author.Email
			}
			info := authorStats[author]
			info.Count++
			if commit.CommittedDate.Time.After(info.LastCommit) {
				info.LastCommit = commit.CommittedDate.Time
			}
			authorStats[author] = info
		}
	}

	return authorStats, nil
}

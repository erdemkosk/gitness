package providers

import (
	"context"
	"fmt"

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

func (g *GitHubProvider) FetchCommits(owner, repo string) (map[string]int, error) {
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
							Nodes []struct {
								Author struct {
									Name  string
									Email string
								}
							}
						} `graphql:"history(first: $limit, after: $after)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"limit": githubv4.Int(100),
		"after": (*githubv4.String)(nil),
	}

	authorStats := make(map[string]int)
	hasNextPage := true

	for hasNextPage {
		err := g.client.Query(context.Background(), &query, variables)
		if err != nil {
			return nil, fmt.Errorf("GitHub query failed: %v", err)
		}

		for _, commit := range query.Repository.DefaultBranchRef.Target.Commit.History.Nodes {
			author := commit.Author.Name
			if author == "" {
				author = commit.Author.Email
			}
			authorStats[author]++
		}

		hasNextPage = query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage
		if hasNextPage {
			variables["after"] = githubv4.String(query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
		}
	}

	return authorStats, nil
}

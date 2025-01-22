package providers

import (
	"context"
	"fmt"

	"github.com/erdemkosk/gitness/internal/constants"
	"github.com/erdemkosk/gitness/internal/progress"
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

func (g *GitHubProvider) FetchCommits(owner, repo string, duration string, branch string) (map[string]CommitInfo, error) {
	var query struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commit struct {
						History struct {
							TotalCount int
							PageInfo   struct {
								HasNextPage bool
								EndCursor   string
							}
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
			} `graphql:"ref(qualifiedName: $ref)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	refName := "refs/heads/master" // default branch
	if branch != "" {
		refName = fmt.Sprintf("refs/heads/%s", branch)
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"limit": githubv4.Int(constants.MaxPageSize),
		"after": (*githubv4.String)(nil),
		"ref":   githubv4.String(refName),
		"since": (*githubv4.GitTimestamp)(nil), // Default olarak nil
	}

	// Add duration if specified
	if duration != "" {
		dur, err := util.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %w", err)
		}
		variables["since"] = githubv4.GitTimestamp{Time: dur.ToTime()}
	}

	authorStats := make(map[string]CommitInfo)
	hasNextPage := true
	processedCommits := 0

	progressBar := progress.NewProgressBar("Fetching commits")
	defer progressBar.Finish()

	// First query to get total count
	err := g.client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, fmt.Errorf("GitHub query failed: %v", err)
	}

	totalCommits := query.Repository.Ref.Target.Commit.History.TotalCount
	progressBar.SetTotal(int64(totalCommits))

	// Process commits from first query
	for _, commit := range query.Repository.Ref.Target.Commit.History.Nodes {
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
		processedCommits++
		progressBar.Increment()
	}

	hasNextPage = query.Repository.Ref.Target.Commit.History.PageInfo.HasNextPage
	if hasNextPage {
		variables["after"] = githubv4.String(query.Repository.Ref.Target.Commit.History.PageInfo.EndCursor)
	}

	for hasNextPage {
		err := g.client.Query(context.Background(), &query, variables)
		if err != nil {
			return nil, fmt.Errorf("GitHub query failed: %v", err)
		}

		for _, commit := range query.Repository.Ref.Target.Commit.History.Nodes {
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
			processedCommits++
			progressBar.Increment()
		}

		hasNextPage = query.Repository.Ref.Target.Commit.History.PageInfo.HasNextPage
		if hasNextPage {
			variables["after"] = githubv4.String(query.Repository.Ref.Target.Commit.History.PageInfo.EndCursor)
		}
	}

	return authorStats, nil
}

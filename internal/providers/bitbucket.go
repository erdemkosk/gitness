package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type BitbucketCommit struct {
	Author struct {
		Raw string `json:"raw"`
	} `json:"author"`
	Date string `json:"date"` // Bitbucket API'den commit tarihi
}

type BitbucketResponse struct {
	Values []BitbucketCommit `json:"values"`
	Next   string            `json:"next"`
}

type BitbucketProvider struct {
	username string
	password string
	client   *http.Client
}

func NewBitbucketProvider(username, password string) *BitbucketProvider {
	return &BitbucketProvider{
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

func (b *BitbucketProvider) FetchCommits(owner, repo string) (map[string]CommitInfo, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo cannot be empty")
	}
	ctx := context.Background()
	authorStats := make(map[string]CommitInfo)
	pageUrl := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commits", owner, repo)

	for pageUrl != "" {
		req, err := http.NewRequestWithContext(ctx, "GET", pageUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		req.SetBasicAuth(b.username, b.password)

		resp, err := b.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commits: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
		}

		var response BitbucketResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %v", err)
		}

		for _, commit := range response.Values {
			author := commit.Author.Raw
			// Bitbucket author format: "Name <email@example.com>"
			if name := strings.Split(author, " <"); len(name) > 0 {
				author = name[0]
			}

			commitDate, err := time.Parse(time.RFC3339, commit.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse commit date: %v", err)
			}

			info := authorStats[author]
			info.Count++
			if commitDate.After(info.LastCommit) {
				info.LastCommit = commitDate
			}
			authorStats[author] = info
		}

		pageUrl = response.Next
	}

	return authorStats, nil
}

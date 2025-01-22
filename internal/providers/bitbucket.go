package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/erdemkosk/gitness/internal/util"
)

type BitbucketProvider struct {
	token  string
	client *http.Client
}

type BitbucketCommit struct {
	Hash   string `json:"hash"`
	Author struct {
		Raw  string `json:"raw"`
		User struct {
			DisplayName string `json:"display_name"`
			AccountID   string `json:"account_id"`
		} `json:"user"`
	} `json:"author"`
	Date string `json:"date"`
}

type BitbucketResponse struct {
	Values []BitbucketCommit `json:"values"`
	Next   string            `json:"next"`
}

func NewBitbucketProvider(clientID, clientSecret string) (*BitbucketProvider, error) {
	token, err := getOAuthToken(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth token: %v", err)
	}

	return &BitbucketProvider{
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func getOAuthToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token, status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (b *BitbucketProvider) FetchCommits(owner, repo string, duration string, branch string) (map[string]CommitInfo, error) {
	authorStats := make(map[string]CommitInfo)

	// Base URL with branch
	var pageURL string
	if branch == "" {
		// For default branch, use main branch endpoint
		pageURL = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commits/refs/heads/master", owner, repo)
	} else {
		// For specific branch
		pageURL = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commits/refs/heads/%s", owner, repo, branch)
	}

	// Query parameters
	params := []string{"pagelen=100"}

	// Add duration parameter if specified
	var since *time.Time
	if duration != "" {
		dur, err := util.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %w", err)
		}
		t := dur.ToTime()
		since = &t
		params = append(params, fmt.Sprintf("since=%s", t.Format(time.RFC3339)))
	}

	// Construct initial URL with parameters
	if len(params) > 0 {
		pageURL += "?" + strings.Join(params, "&")
	}

	for pageURL != "" {
		commits, nextPage, err := b.fetchPage(pageURL)
		if err != nil {
			// If master branch fails, try main branch
			if branch == "" && strings.Contains(pageURL, "refs/heads/master") {
				pageURL = strings.Replace(pageURL, "refs/heads/master", "refs/heads/main", 1)
				commits, nextPage, err = b.fetchPage(pageURL)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch commits from both master and main branches: %v", err)
				}
			} else {
				return nil, fmt.Errorf("failed to fetch commits: %v", err)
			}
		}

		for _, commit := range commits {
			commitDate, err := time.Parse(time.RFC3339, commit.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse commit date: %v", err)
			}

			if since != nil && commitDate.Before(*since) {
				continue
			}

			author := b.extractAuthorName(commit)
			info := authorStats[author]
			info.Count++
			if commitDate.After(info.LastCommit) {
				info.LastCommit = commitDate
			}
			authorStats[author] = info
		}

		pageURL = nextPage
	}

	return authorStats, nil
}

func (b *BitbucketProvider) fetchPage(url string) ([]BitbucketCommit, string, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("authentication error (401). API response: %s", string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API request failed, status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var response BitbucketResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return response.Values, response.Next, nil
}

func (b *BitbucketProvider) extractAuthorName(commit BitbucketCommit) string {
	if commit.Author.User.DisplayName != "" {
		return commit.Author.User.DisplayName
	}

	if commit.Author.Raw != "" {
		parts := strings.Split(commit.Author.Raw, " <")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if commit.Author.User.AccountID != "" {
		return commit.Author.User.AccountID
	}

	return "Unknown Author"
}

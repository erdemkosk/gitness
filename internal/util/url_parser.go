package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Duration struct {
	Value int
	Unit  string // "day", "month", "year"
}

type RepoInfo struct {
	Owner        string
	Repo         string
	ProviderType string
	Config       map[string]string
	Duration     *Duration
}

func ParseRepositoryURL(url string) (*RepoInfo, error) {
	url = strings.TrimSpace(url)

	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository URL: %s", url)
	}

	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)

	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner or repo cannot be empty. URL: %s", url)
	}

	info := &RepoInfo{
		Owner: owner,
		Repo:  repo,
	}

	if strings.Contains(url, "github.com") {
		info.ProviderType = "github"
		info.Config = map[string]string{
			"token": os.Getenv("GITHUB_TOKEN"),
		}
	} else if strings.Contains(url, "bitbucket.org") {
		clientID := os.Getenv("BITBUCKET_CLIENT_ID")
		clientSecret := os.Getenv("BITBUCKET_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			return nil, fmt.Errorf("BITBUCKET_CLIENT_ID and BITBUCKET_CLIENT_SECRET environment variables are required for Bitbucket")
		}

		info.ProviderType = "bitbucket"
		info.Config = map[string]string{
			"clientID":     clientID,
			"clientSecret": clientSecret,
		}
	} else {
		return nil, fmt.Errorf("unsupported repository provider: %s", url)
	}

	return info, nil
}

func ParseDuration(duration string) (*Duration, error) {
	if duration == "" {
		return nil, nil
	}

	var value int
	var unit string
	_, err := fmt.Sscanf(duration, "%d%s", &value, &unit)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format. Use: 6m, 1y, 30d")
	}

	switch unit {
	case "d", "day", "days":
		return &Duration{Value: value, Unit: "day"}, nil
	case "m", "month", "months":
		return &Duration{Value: value, Unit: "month"}, nil
	case "y", "year", "years":
		return &Duration{Value: value, Unit: "year"}, nil
	default:
		return nil, fmt.Errorf("invalid duration unit. Use: d(days), m(months), y(years)")
	}
}

func (d *Duration) ToTime() time.Time {
	now := time.Now()
	switch d.Unit {
	case "day":
		return now.AddDate(0, 0, -d.Value)
	case "month":
		return now.AddDate(0, -d.Value, 0)
	case "year":
		return now.AddDate(-d.Value, 0, 0)
	}
	return now
}

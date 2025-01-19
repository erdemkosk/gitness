package util

import (
	"fmt"
	"os"
	"strings"
)

type RepoInfo struct {
	Owner        string
	Repo         string
	ProviderType string
	Config       map[string]string
}

func ParseRepositoryURL(url string) (*RepoInfo, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository URL")
	}

	info := &RepoInfo{
		Owner: parts[len(parts)-2],
		Repo:  parts[len(parts)-1],
	}

	if strings.Contains(url, "github.com") {
		info.ProviderType = "github"
		info.Config = map[string]string{
			"token": os.Getenv("GITHUB_TOKEN"),
		}
	} else if strings.Contains(url, "bitbucket.org") {
		info.ProviderType = "bitbucket"
		info.Config = map[string]string{
			"username": os.Getenv("BITBUCKET_USERNAME"),
			"password": os.Getenv("BITBUCKET_PASSWORD"),
		}
	} else {
		return nil, fmt.Errorf("unsupported repository host")
	}

	return info, nil
}

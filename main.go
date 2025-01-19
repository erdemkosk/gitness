package main

import (
	"log"
	"os"
	"strings"

	"github.com/erdemkosk/gitness/internal/analyzer"
	"github.com/erdemkosk/gitness/internal/providers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) != 2 {
		log.Fatal("Please provide repository URL as an argument")
	}

	url := os.Args[1]
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		log.Fatal("Invalid repository URL")
	}

	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Determine provider based on URL
	var provider providers.CommitProvider
	if strings.Contains(url, "github.com") {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			log.Fatal("GITHUB_TOKEN environment variable is required")
		}
		provider = providers.NewGitHubProvider(token)
	} else if strings.Contains(url, "bitbucket.org") {
		username := os.Getenv("BITBUCKET_USERNAME")
		password := os.Getenv("BITBUCKET_PASSWORD")
		if username == "" || password == "" {
			log.Fatal("BITBUCKET_USERNAME and BITBUCKET_PASSWORD environment variables are required")
		}
		provider = providers.NewBitbucketProvider(username, password)
	} else {
		log.Fatal("Unsupported repository host")
	}

	repoAnalyzer := analyzer.NewRepositoryAnalyzer(provider)
	stats, err := repoAnalyzer.Analyze(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	stats.Print()
}

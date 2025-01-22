package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OutputFormat string
	Duration     string
	Branch       string
	RepoURL      string
}

func LoadConfig() (*Config, error) {
	outputFormat := flag.String("output", "", "Output format: console, json, or markdown")
	duration := flag.String("duration", "", "Analyze commits for last duration (e.g., 6m, 1y, 30d)")
	branch := flag.String("branch", "", "Specific branch to analyze")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: Error loading .env file: %v\n", err)
		}
	}

	config := &Config{}

	// Output Format
	config.OutputFormat = getOutputFormat(*outputFormat)

	// Duration
	config.Duration = getDuration(*duration)

	// Branch
	config.Branch = getBranch(*branch)

	// Repository URL
	repoURL, err := getRepositoryURL()
	if err != nil {
		return nil, err
	}
	config.RepoURL = repoURL

	return config, nil
}

func getOutputFormat(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envFormat := os.Getenv("OUTPUT_FORMAT"); envFormat != "" {
		return envFormat
	}
	return "console"
}

func getDuration(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("COMMIT_HISTORY_DURATION")
}

func getBranch(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("REPOSITORY_BRANCH")
}

func getRepositoryURL() (string, error) {
	args := flag.Args()
	if len(args) == 1 {
		return args[0], nil
	}

	if url := os.Getenv("REPOSITORY_URL"); url != "" {
		return url, nil
	}

	return "", fmt.Errorf("please provide repository URL either as argument or set REPOSITORY_URL environment variable")
}

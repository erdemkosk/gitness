package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/erdemkosk/gitness/internal/analyzer"
	"github.com/erdemkosk/gitness/internal/output"
	"github.com/erdemkosk/gitness/internal/providers"
	"github.com/erdemkosk/gitness/internal/util"
	"github.com/joho/godotenv"
)

func getRepositoryURL() (string, error) {
	// First try command line argument
	args := flag.Args()
	if len(args) == 1 {
		return args[0], nil
	}

	// If no args, try environment variable
	if url := os.Getenv("REPOSITORY_URL"); url != "" {
		return url, nil
	}

	return "", fmt.Errorf("please provide repository URL either as argument or set REPOSITORY_URL environment variable")
}

func main() {
	// Command line flag has higher priority
	outputFormat := flag.String("output", "", "Output format: console, json, or markdown")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		// Ignore .env file error in Docker environment
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	// If output format is not provided via flag, check environment variable
	if *outputFormat == "" {
		envFormat := os.Getenv("OUTPUT_FORMAT")
		if envFormat != "" {
			*outputFormat = envFormat
		} else {
			// Default to console if not specified anywhere
			*outputFormat = "console"
		}
	}

	url, err := getRepositoryURL()
	if err != nil {
		log.Fatal(err)
	}

	repoInfo, err := util.ParseRepositoryURL(url)
	if err != nil {
		log.Fatal(err)
	}

	providerFactory := providers.NewProviderFactory()
	provider, err := providerFactory.CreateProvider(repoInfo.ProviderType, repoInfo.Config)
	if err != nil {
		log.Fatal(err)
	}

	repoAnalyzer := analyzer.NewRepositoryAnalyzer(provider)
	stats, err := repoAnalyzer.Analyze(repoInfo.Owner, repoInfo.Repo)
	if err != nil {
		log.Fatal(err)
	}

	formatterFactory := output.NewFormatterFactory()
	formatter, exists := formatterFactory.GetFormatter(*outputFormat)
	if !exists {
		log.Fatalf("Unknown output format: %s", *outputFormat)
	}

	result, err := formatter.Format(stats)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

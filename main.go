package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/erdemkosk/gitness/internal/analyzer"
	"github.com/erdemkosk/gitness/internal/output"
	"github.com/erdemkosk/gitness/internal/providers"
	"github.com/erdemkosk/gitness/internal/util"
	"github.com/joho/godotenv"
)

func main() {
	outputFormat := flag.String("output", "console", "Output format: console, json, or markdown")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Please provide repository URL as an argument")
	}
	url := args[0]

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

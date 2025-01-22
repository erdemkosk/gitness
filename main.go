package main

import (
	"fmt"
	"log"

	"github.com/erdemkosk/gitness/internal/analyzer"
	"github.com/erdemkosk/gitness/internal/config"
	"github.com/erdemkosk/gitness/internal/output"
	"github.com/erdemkosk/gitness/internal/providers"
	"github.com/erdemkosk/gitness/internal/util"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	repoInfo, err := util.ParseRepositoryURL(cfg.RepoURL)
	if err != nil {
		log.Fatal(err)
	}

	providerFactory := providers.NewProviderFactory()
	provider, err := providerFactory.CreateProvider(repoInfo.ProviderType, repoInfo.Config)
	if err != nil {
		log.Fatal(err)
	}

	repoAnalyzer := analyzer.NewRepositoryAnalyzer(provider)
	stats, err := repoAnalyzer.Analyze(repoInfo.Owner, repoInfo.Repo, cfg.Duration, cfg.Branch)
	if err != nil {
		log.Fatal(err)
	}

	formatterFactory := output.NewFormatterFactory()
	formatter, exists := formatterFactory.GetFormatter(cfg.OutputFormat)
	if !exists {
		log.Fatalf("Unknown output format: %s", cfg.OutputFormat)
	}

	result, err := formatter.Format(stats)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

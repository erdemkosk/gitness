package analyzer

import (
	"fmt"
	"sort"

	"github.com/erdemkosk/gitness/internal/models"
	"github.com/erdemkosk/gitness/internal/providers"
)

type RepositoryAnalyzer struct {
	provider providers.CommitProvider
}

func NewRepositoryAnalyzer(provider providers.CommitProvider) *RepositoryAnalyzer {
	if provider == nil {
		panic("provider cannot be nil")
	}
	return &RepositoryAnalyzer{
		provider: provider,
	}
}

func (ra *RepositoryAnalyzer) Analyze(owner, repo string) (*models.RepositoryStats, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo cannot be empty")
	}

	stats, err := ra.provider.FetchCommits(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}

	if len(stats) == 0 {
		return nil, fmt.Errorf("no commits found in repository")
	}

	total := 0
	for _, count := range stats {
		total += count
	}

	var contributors []models.Contributor
	for author, count := range stats {
		percentage := float64(count) * 100 / float64(total)
		contributors = append(contributors, models.Contributor{
			Name:       author,
			Commits:    count,
			Percentage: percentage,
		})
	}

	// Sort contributors by commit count in descending order
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})

	return &models.RepositoryStats{
		Owner:        owner,
		Repo:         repo,
		Contributors: contributors,
		BusFactor:    calculateBusFactor(stats),
		TotalCommits: total,
	}, nil
}

func calculateBusFactor(stats map[string]int) int {
	if len(stats) == 0 {
		return 0
	}

	var total int
	for _, count := range stats {
		total += count
	}

	type authorStat struct {
		name       string
		percentage float64
	}

	var statsList []authorStat
	for author, count := range stats {
		percentage := float64(count) * 100 / float64(total)
		statsList = append(statsList, authorStat{author, percentage})
	}

	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].percentage > statsList[j].percentage
	})

	var cumulative float64
	var busFactor int
	for _, stat := range statsList {
		cumulative += stat.percentage
		busFactor++
		if cumulative >= 80 {
			break
		}
	}

	return busFactor
}

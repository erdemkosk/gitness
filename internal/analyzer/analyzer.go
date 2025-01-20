package analyzer

import (
	"fmt"
	"math"
	"sort"
	"time"

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

func (ra *RepositoryAnalyzer) Analyze(owner, repo string, duration string) (*models.RepositoryStats, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo cannot be empty")
	}

	stats, err := ra.provider.FetchCommits(owner, repo, duration)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}

	if len(stats) == 0 {
		return nil, fmt.Errorf("no commits found in repository")
	}

	total := 0
	for _, info := range stats {
		total += info.Count
	}

	var contributors []models.Contributor
	for author, info := range stats {
		percentage := float64(info.Count) * 100 / float64(total)
		contributors = append(contributors, models.Contributor{
			Name:       author,
			Commits:    info.Count,
			Percentage: percentage,
			LastCommit: info.LastCommit,
		})
	}

	// Sort contributors by commit count in descending order
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})

	activeContributors := 0
	recentContributors := 0
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	for _, contributor := range contributors {
		if contributor.Percentage >= 1.0 {
			activeContributors++
		}
		if contributor.LastCommit.After(threeMonthsAgo) {
			recentContributors++
		}
	}

	contributorActivity := float64(activeContributors) / float64(len(contributors)) * 100

	// Knowledge Distribution Score hesaplama
	knowledgeScore := calculateKnowledgeDistribution(stats)

	return &models.RepositoryStats{
		Owner:               owner,
		Repo:                repo,
		Contributors:        contributors,
		BusFactor:           calculateBusFactor(stats),
		TotalCommits:        total,
		ContributorActivity: contributorActivity,
		RecentContributors:  recentContributors,
		KnowledgeScore:      knowledgeScore,
		AnalysisDuration:    duration,
	}, nil
}

func calculateBusFactor(stats map[string]providers.CommitInfo) int {
	if len(stats) == 0 {
		return 0
	}

	var total int
	for _, info := range stats {
		total += info.Count
	}

	type authorStat struct {
		name       string
		percentage float64
	}

	var statsList []authorStat
	for author, info := range stats {
		percentage := float64(info.Count) * 100 / float64(total)
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

func calculateKnowledgeDistribution(stats map[string]providers.CommitInfo) float64 {
	if len(stats) == 0 {
		return 0
	}

	var total int
	for _, info := range stats {
		total += info.Count
	}

	var sumDifferences float64
	var n = float64(len(stats))

	for _, info1 := range stats {
		percent1 := float64(info1.Count) / float64(total)
		for _, info2 := range stats {
			percent2 := float64(info2.Count) / float64(total)
			sumDifferences += math.Abs(percent1 - percent2)
		}
	}

	gini := sumDifferences / (2 * n * n)
	knowledgeScore := (1 - gini) * 100

	return math.Round(knowledgeScore*100) / 100
}

package analyzer

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/erdemkosk/gitness/internal/constants"
	"github.com/erdemkosk/gitness/internal/models"
	"github.com/erdemkosk/gitness/internal/providers"
	"github.com/erdemkosk/gitness/internal/util"
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

func (ra *RepositoryAnalyzer) analyzeCommitFrequency(contributors []models.Contributor, totalCommits int, duration string) (float64, float64, float64, string, string) {
	if len(contributors) == 0 {
		return 0, 0, 0, "", ""
	}

	// Statistics by day and hour
	dayCount := make(map[string]int)
	hourCount := make(map[int]int)

	// Find oldest and newest commit
	var firstCommit, lastCommit time.Time
	firstCommit = time.Now()
	lastCommit = time.Time{}

	for _, contributor := range contributors {
		commitTime := contributor.LastCommit

		// Update oldest and newest commit times
		if commitTime.Before(firstCommit) {
			firstCommit = commitTime
		}
		if commitTime.After(lastCommit) {
			lastCommit = commitTime
		}

		// Day and hour statistics
		dayCount[commitTime.Weekday().String()]++
		hourCount[commitTime.Hour()]++
	}

	// Calculate time period
	var timePeriod time.Duration
	if duration != "" {
		// Calculate based on provided duration
		dur, err := util.ParseDuration(duration)
		if err == nil {
			timePeriod = time.Since(dur.ToTime())
		}
	}

	if timePeriod == 0 {
		// If no duration or invalid, use time between first and last commit
		timePeriod = lastCommit.Sub(firstCommit)
	}

	// Calculate days, weeks, and months
	days := math.Max(math.Ceil(timePeriod.Hours()/24), 1)
	weeks := math.Max(math.Ceil(days/7), 1)
	months := math.Max(math.Ceil(days/30), 1)

	// Calculate averages
	dailyAvg := float64(totalCommits) / days
	weeklyAvg := float64(totalCommits) / weeks
	monthlyAvg := float64(totalCommits) / months

	// Round to 2 decimal places
	dailyAvg = math.Round(dailyAvg*100) / 100
	weeklyAvg = math.Round(weeklyAvg*100) / 100
	monthlyAvg = math.Round(monthlyAvg*100) / 100

	// Find most active day and hour
	mostActiveDay := ""
	maxDayCount := 0
	for day, count := range dayCount {
		if count > maxDayCount {
			maxDayCount = count
			mostActiveDay = day
		}
	}

	mostActiveHour := 0
	maxHourCount := 0
	for hour, count := range hourCount {
		if count > maxHourCount {
			maxHourCount = count
			mostActiveHour = hour
		}
	}

	mostActiveTime := fmt.Sprintf("%02d:00", mostActiveHour)

	return dailyAvg, weeklyAvg, monthlyAvg, mostActiveDay, mostActiveTime
}

func (ra *RepositoryAnalyzer) Analyze(owner, repo string, duration string, branch string) (*models.RepositoryStats, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("%s", constants.ErrorMessages.EmptyOwnerRepo)
	}

	commits, err := ra.provider.FetchCommits(owner, repo, duration, branch)
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("%s", constants.ErrorMessages.NoCommitsFound)
	}

	total := 0
	for _, info := range commits {
		total += info.Count
	}

	var contributors []models.Contributor
	for author, info := range commits {
		percentage := float64(info.Count) * 100 / float64(total)
		contributors = append(contributors, models.Contributor{
			Name:       author,
			Commits:    info.Count,
			Percentage: percentage,
			LastCommit: info.LastCommit,
		})
	}

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})

	dailyAvg, weeklyAvg, monthlyAvg, mostActiveDay, mostActiveTime := ra.analyzeCommitFrequency(contributors, total, duration)

	activeContributors := 0
	recentContributors := 0
	threeMonthsAgo := time.Now().AddDate(0, -constants.Metrics.ContributorStats.RecentMonths, 0)

	for _, contributor := range contributors {
		if contributor.Percentage > constants.Metrics.ContributorStats.SignificantContrib {
			activeContributors++
		}
		if contributor.LastCommit.After(threeMonthsAgo) {
			recentContributors++
		}
	}

	contributorActivity := float64(activeContributors) / float64(len(contributors)) * constants.Metrics.KnowledgeScore.Maximum
	knowledgeScore := calculateKnowledgeDistribution(commits)

	return &models.RepositoryStats{
		Owner:                owner,
		Repo:                 repo,
		Branch:               branch,
		Contributors:         contributors,
		BusFactor:            calculateBusFactor(commits),
		TotalCommits:         total,
		ContributorActivity:  contributorActivity,
		RecentContributors:   recentContributors,
		KnowledgeScore:       knowledgeScore,
		AnalysisDuration:     duration,
		DailyCommitAverage:   math.Round(dailyAvg*100) / 100,
		WeeklyCommitAverage:  math.Round(weeklyAvg*100) / 100,
		MonthlyCommitAverage: math.Round(monthlyAvg*100) / 100,
		MostActiveDay:        mostActiveDay,
		MostActiveTime:       mostActiveTime,
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
		if cumulative >= constants.Metrics.BusFactor.Percentage {
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
	knowledgeScore := (1 - gini) * constants.Metrics.KnowledgeScore.Maximum

	return math.Round(knowledgeScore*100) / 100
}

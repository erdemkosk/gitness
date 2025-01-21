package output

import (
	"encoding/json"

	"github.com/erdemkosk/gitness/internal/models"
)

type JSONFormatter struct{}

type jsonOutput struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Analysis struct {
		Duration            string  `json:"duration"`
		BusFactor           int     `json:"busFactor"`
		KnowledgeScore      float64 `json:"knowledgeScore"`
		ContributorActivity float64 `json:"contributorActivity"`
		RecentContributors  int     `json:"recentContributors"`
		TotalCommits        int     `json:"totalCommits"`
	} `json:"analysis"`
	CommitFrequency struct {
		DailyAverage   float64 `json:"dailyAverage"`
		WeeklyAverage  float64 `json:"weeklyAverage"`
		MonthlyAverage float64 `json:"monthlyAverage"`
		MostActiveDay  string  `json:"mostActiveDay"`
		PeakTime       string  `json:"peakTime"`
	} `json:"commitFrequency"`
	Contributors []struct {
		Name       string  `json:"name"`
		Commits    int     `json:"commits"`
		Percentage float64 `json:"percentage"`
		LastCommit string  `json:"lastCommit"`
	} `json:"contributors"`
}

func (f *JSONFormatter) Format(stats *models.RepositoryStats) (string, error) {
	output := jsonOutput{
		Owner: stats.Owner,
		Repo:  stats.Repo,
	}

	output.Analysis.Duration = stats.AnalysisDuration
	output.Analysis.BusFactor = stats.BusFactor
	output.Analysis.KnowledgeScore = stats.KnowledgeScore
	output.Analysis.ContributorActivity = stats.ContributorActivity
	output.Analysis.RecentContributors = stats.RecentContributors
	output.Analysis.TotalCommits = stats.TotalCommits

	output.CommitFrequency.DailyAverage = stats.DailyCommitAverage
	output.CommitFrequency.WeeklyAverage = stats.WeeklyCommitAverage
	output.CommitFrequency.MonthlyAverage = stats.MonthlyCommitAverage
	output.CommitFrequency.MostActiveDay = stats.MostActiveDay
	output.CommitFrequency.PeakTime = stats.MostActiveTime

	output.Contributors = make([]struct {
		Name       string  `json:"name"`
		Commits    int     `json:"commits"`
		Percentage float64 `json:"percentage"`
		LastCommit string  `json:"lastCommit"`
	}, len(stats.Contributors))

	for i, c := range stats.Contributors {
		output.Contributors[i] = struct {
			Name       string  `json:"name"`
			Commits    int     `json:"commits"`
			Percentage float64 `json:"percentage"`
			LastCommit string  `json:"lastCommit"`
		}{
			Name:       c.Name,
			Commits:    c.Commits,
			Percentage: c.Percentage,
			LastCommit: c.LastCommit.Format("2006-01-02"),
		}
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

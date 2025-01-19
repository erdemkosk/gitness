package output

import (
	"encoding/json"

	"github.com/erdemkosk/gitness/internal/models"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(stats *models.RepositoryStats) (string, error) {
	data := struct {
		Owner               string               `json:"owner"`
		Repo                string               `json:"repo"`
		BusFactor           int                  `json:"busFactor"`
		TotalCommits        int                  `json:"totalCommits"`
		Contributors        []models.Contributor `json:"contributors"`
		ContributorActivity float64              `json:"contributorActivity"` // Percentage of contributors with >1% contribution
		RecentContributors  int                  `json:"recentContributors"`  // Number of contributors active in last 3 months
	}{
		Owner:               stats.Owner,
		Repo:                stats.Repo,
		BusFactor:           stats.BusFactor,
		TotalCommits:        stats.TotalCommits,
		Contributors:        stats.Contributors,
		ContributorActivity: stats.ContributorActivity,
		RecentContributors:  stats.RecentContributors,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

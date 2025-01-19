package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Author struct {
	Name  string
	Email string
}

type Commit struct {
	Author Author
}

type Contributor struct {
	Name       string
	Commits    int
	Percentage float64
	LastCommit time.Time
}

type RepositoryStats struct {
	Owner               string
	Repo                string
	Contributors        []Contributor
	BusFactor           int
	TotalCommits        int
	ContributorActivity float64
	RecentContributors  int
}

func (rs *RepositoryStats) Print() {
	fmt.Printf("\nRepo: %s/%s\n", rs.Owner, rs.Repo)
	fmt.Printf("Bus Factor: %d\n\n", rs.BusFactor)
	fmt.Printf("Contributors:\n")
	fmt.Printf("------------------\n")

	for _, c := range rs.Contributors {
		fmt.Printf("%s: %d commits (%.1f%%)\n", c.Name, c.Commits, c.Percentage)
	}
}

func (rs *RepositoryStats) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Owner        string        `json:"owner"`
		Repo         string        `json:"repo"`
		BusFactor    int           `json:"busFactor"`
		TotalCommits int           `json:"totalCommits"`
		Contributors []Contributor `json:"contributors"`
	}{
		Owner:        rs.Owner,
		Repo:         rs.Repo,
		BusFactor:    rs.BusFactor,
		TotalCommits: rs.TotalCommits,
		Contributors: rs.Contributors,
	})
}

func (rs *RepositoryStats) ToMarkdown() string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# Repository Analysis: %s/%s\n\n", rs.Owner, rs.Repo))
	md.WriteString(fmt.Sprintf("## Bus Factor: %d\n\n", rs.BusFactor))
	md.WriteString(fmt.Sprintf("Total Commits: %d\n\n", rs.TotalCommits))

	md.WriteString("## Contributors\n\n")
	md.WriteString("| Name | Commits | Percentage |\n")
	md.WriteString("|------|---------|------------|\n")

	for _, c := range rs.Contributors {
		md.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", c.Name, c.Commits, c.Percentage))
	}

	return md.String()
}

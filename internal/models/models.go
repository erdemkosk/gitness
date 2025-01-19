package models

import "fmt"

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
}

type RepositoryStats struct {
	Owner        string
	Repo         string
	Contributors []Contributor
	BusFactor    int
	TotalCommits int
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

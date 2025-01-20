package output

import (
	"fmt"
	"strings"

	"github.com/erdemkosk/gitness/internal/models"
)

type ConsoleFormatter struct{}

func (f *ConsoleFormatter) Format(stats *models.RepositoryStats) (string, error) {
	var output strings.Builder

	output.WriteString(` 
 :................:++-......................................................  
 :..................:...#%..................................................  
 :........=%@@#+%-.*#.=%@@%%%%%%%%--%+*@@@+:...-#@@%=..:+%@@*:.:#@@%=.......  
 :.......%#:..:%@-.#%...#%....::::.-@%:...*@:.%%:...#%:*@:..:..@*...:.......  
 :......-@:....-@-.#%...#%......::.-@=....:@=-@@@@@@@%-:#%@@+..-#@@%=.......  
 :.......%#:..:%@-.#%...*%:.....::.-@=....:@=:@#:...:..::...#%.:...-@-......  
 :........-%%%#+@-.*%....*@@%..+++--@=....:@=..+@@@@%-.=%@@@%:.*@@@@#.......  
 :.......-+:::-@#...........................................................  
 :........-*##+:............................................................    
`)

	output.WriteString(fmt.Sprintf("\nRepo: %s/%s\n", stats.Owner, stats.Repo))

	if stats.AnalysisDuration != "" {
		output.WriteString(fmt.Sprintf("\nAnalysis Period: Last %s\n", stats.AnalysisDuration))
	} else {
		output.WriteString("\nAnalysis Period: All Time\n")
	}

	// ANSI color codes
	red := "\033[31m"
	yellow := "\033[33m"
	green := "\033[32m"
	reset := "\033[0m"

	// Bus Factor coloring
	busFactorColor := green
	if stats.BusFactor < 2 {
		busFactorColor = red
	} else if stats.BusFactor < 4 {
		busFactorColor = yellow
	}

	// Active Contributor Ratio coloring
	activityColor := green
	if stats.ContributorActivity < 30 {
		activityColor = red
	} else if stats.ContributorActivity < 50 {
		activityColor = yellow
	}

	// Recent Contributors coloring
	recentColor := green
	if stats.RecentContributors < 2 {
		recentColor = red
	} else if stats.RecentContributors < 5 {
		recentColor = yellow
	}

	output.WriteString(fmt.Sprintf("Bus Factor: %s%d%s (critical if < 2, warning if < 4)\n",
		busFactorColor, stats.BusFactor, reset))

	output.WriteString(fmt.Sprintf("Knowledge Distribution Score: %s%.1f%s (0-100, higher is better)\n",
		getScoreColor(stats.KnowledgeScore), stats.KnowledgeScore, reset))

	output.WriteString(fmt.Sprintf("Active Contributor Ratio: %s%.1f%%%s (contributors with >1%% contribution, critical if < 30%%, warning if < 50%%)\n",
		activityColor, stats.ContributorActivity, reset))

	output.WriteString(fmt.Sprintf("Recent Contributors: %s%d%s (active in last 3 months, critical if < 2, warning if < 5)\n\n",
		recentColor, stats.RecentContributors, reset))

	output.WriteString("Contributors:\n")
	output.WriteString("------------------\n")

	for _, c := range stats.Contributors {
		output.WriteString(fmt.Sprintf("%s: %d commits (%.1f%%)\n",
			c.Name, c.Commits, c.Percentage))
	}

	return output.String(), nil
}

func getScoreColor(score float64) string {
	if score < 25 {
		return "\033[31m" // red
	} else if score < 50 {
		return "\033[33m" // yellow
	}
	return "\033[32m" // green
}

package output

import (
	"fmt"
	"strings"

	"github.com/erdemkosk/gitness/internal/models"
)

type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(stats *models.RepositoryStats) (string, error) {
	var md strings.Builder

	md.WriteString("![Gitness](https://github.com/erdemkosk/gitness/blob/master/logo.png?raw=true)\n\n")
	md.WriteString(fmt.Sprintf("# Repository Analysis: %s/%s\n\n", stats.Owner, stats.Repo))

	if stats.AnalysisDuration != "" {
		md.WriteString(fmt.Sprintf("## 📅 Analysis Period: Last %s\n\n", stats.AnalysisDuration))
	} else {
		md.WriteString("## 📅 Analysis Period: All Time\n\n")
	}

	// Bus Factor status
	busFactorEmoji := "🟢"
	if stats.BusFactor < 2 {
		busFactorEmoji = "🔴"
	} else if stats.BusFactor < 4 {
		busFactorEmoji = "🟡"
	}

	// Active Contributor status
	activityEmoji := "🟢"
	if stats.ContributorActivity < 30 {
		activityEmoji = "🔴"
	} else if stats.ContributorActivity < 50 {
		activityEmoji = "🟡"
	}

	// Recent Contributors status
	recentEmoji := "🟢"
	if stats.RecentContributors < 2 {
		recentEmoji = "🔴"
	} else if stats.RecentContributors < 5 {
		recentEmoji = "🟡"
	}

	// Knowledge Distribution status
	scoreEmoji := "🟢"
	if stats.KnowledgeScore < 25 {
		scoreEmoji = "🔴"
	} else if stats.KnowledgeScore < 50 {
		scoreEmoji = "🟡"
	}

	md.WriteString(fmt.Sprintf("## %s Bus Factor: **%d** (critical if < 2, warning if < 4)\n\n",
		busFactorEmoji, stats.BusFactor))

	md.WriteString(fmt.Sprintf("## %s Knowledge Distribution Score: **%.1f** (0-100, higher is better)\n\n",
		scoreEmoji, stats.KnowledgeScore))

	md.WriteString(fmt.Sprintf("## %s Active Contributor Ratio: **%.1f%%** (contributors with >1%% contribution, critical if < 30%%, warning if < 50%%)\n\n",
		activityEmoji, stats.ContributorActivity))

	md.WriteString(fmt.Sprintf("## %s Recent Contributors: **%d** (active in last 3 months, critical if < 2, warning if < 5)\n\n",
		recentEmoji, stats.RecentContributors))

	md.WriteString(fmt.Sprintf("Total Commits: %d\n\n", stats.TotalCommits))

	md.WriteString("## Contributors\n\n")
	md.WriteString("| Name | Commits | Percentage |\n")
	md.WriteString("|------|---------|------------|\n")

	for _, c := range stats.Contributors {
		md.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n",
			c.Name, c.Commits, c.Percentage))
	}

	return md.String(), nil
}

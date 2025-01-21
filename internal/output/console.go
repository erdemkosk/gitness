package output

import (
	"fmt"
	"strings"

	"github.com/erdemkosk/gitness/internal/constants"
	"github.com/erdemkosk/gitness/internal/models"
	"github.com/fatih/color"
)

type ConsoleFormatter struct{}

func (f *ConsoleFormatter) Format(stats *models.RepositoryStats) (string, error) {
	var output strings.Builder

	// Colors
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Logo and title
	output.WriteString("\n")
	output.WriteString(blue(`
    ██████╗ ██╗████████╗███╗   ██╗███████╗███████╗███████╗
   ██╔════╝ ██║╚══██╔══╝████╗  ██║██╔════╝██╔════╝██╔════╝
   ██║  ███╗██║   ██║   ██╔██╗ ██║█████╗  ███████╗███████╗
   ██║   ██║██║   ██║   ██║╚██╗██║██╔══╝  ╚════██║╚════██║
   ╚██████╔╝██║   ██║   ██║ ╚████║███████╗███████║███████║
    ╚═════╝ ╚═╝   ╚═╝   ╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝
`))
	output.WriteString("\n")

	// Title
	output.WriteString(magenta("Your repo's fitness witness! Track your bus factor before your code misses the bus."))
	output.WriteString("\n")

	// Repository info
	output.WriteString(yellow("═══════════════════════════════════════════════════════════════\n\n"))
	output.WriteString(cyan("📊 Repository: "))
	output.WriteString(cyan(fmt.Sprintf("%s/%s", stats.Owner, stats.Repo)))
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 50) + "\n")

	if stats.AnalysisDuration != "" {
		output.WriteString(yellow(fmt.Sprintf("🕒 Analysis Period: Last %s\n", stats.AnalysisDuration)))
	} else {
		output.WriteString(yellow("🕒 Analysis Period: All Time\n"))
	}
	output.WriteString("\n")

	// Core Metrics
	output.WriteString(cyan("🎯 Core Metrics\n"))
	output.WriteString(strings.Repeat("─", 20) + "\n")

	busFactorColor := green
	if stats.BusFactor < constants.BusFactorCriticalThreshold {
		busFactorColor = red
	} else if stats.BusFactor < constants.BusFactorWarningThreshold {
		busFactorColor = yellow
	}

	knowledgeScoreColor := green
	if stats.KnowledgeScore < constants.KnowledgeScoreCriticalThreshold {
		knowledgeScoreColor = red
	} else if stats.KnowledgeScore < constants.KnowledgeScoreWarningThreshold {
		knowledgeScoreColor = yellow
	}

	output.WriteString(fmt.Sprintf("🚌 Bus Factor: %s\n", busFactorColor(stats.BusFactor)))
	output.WriteString(fmt.Sprintf("📚 Knowledge Distribution: %s%%\n", knowledgeScoreColor(stats.KnowledgeScore)))
	output.WriteString(fmt.Sprintf("📝 Total Commits: %s\n", yellow(stats.TotalCommits)))
	output.WriteString(fmt.Sprintf("👥 Active Contributors: %s%%\n", yellow(stats.ContributorActivity)))
	output.WriteString(fmt.Sprintf("🔄 Recent Contributors: %s\n", yellow(stats.RecentContributors)))
	output.WriteString("\n")

	// Commit Frequency
	output.WriteString(cyan("⏰ Commit Frequency\n"))
	output.WriteString(strings.Repeat("─", 20) + "\n")
	output.WriteString(fmt.Sprintf("📅 Daily Average: %s commits\n", yellow(fmt.Sprintf("%.2f", stats.DailyCommitAverage))))
	output.WriteString(fmt.Sprintf("📅 Weekly Average: %s commits\n", yellow(fmt.Sprintf("%.2f", stats.WeeklyCommitAverage))))
	output.WriteString(fmt.Sprintf("📅 Monthly Average: %s commits\n", yellow(fmt.Sprintf("%.2f", stats.MonthlyCommitAverage))))
	output.WriteString(fmt.Sprintf("📅 Most Active Day: %s\n", yellow(stats.MostActiveDay)))
	output.WriteString(fmt.Sprintf("🕒 Peak Activity Time: %s\n", yellow(stats.MostActiveTime)))
	output.WriteString("\n")

	// Contributors
	output.WriteString(cyan("👥 Top Contributors\n"))
	output.WriteString(strings.Repeat("─", 20) + "\n")

	for i, c := range stats.Contributors {
		if i >= constants.MaxDisplayedContributors {
			break
		}
		output.WriteString(fmt.Sprintf("👤 %s: %s commits (%s%%)\n",
			yellow(c.Name),
			green(fmt.Sprintf("%d", c.Commits)),
			yellow(fmt.Sprintf("%.1f", c.Percentage))))
	}

	if len(stats.Contributors) > constants.MaxDisplayedContributors {
		output.WriteString(fmt.Sprintf("\n... and %d more contributors\n", len(stats.Contributors)-constants.MaxDisplayedContributors))
	}

	// Footer
	output.WriteString("\n")
	output.WriteString(yellow("═══════════════════════════════════════════════════════════════\n"))
	output.WriteString(magenta("                     Generated by Gitness                      \n"))

	return output.String(), nil
}

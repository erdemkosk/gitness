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
	md.WriteString(fmt.Sprintf("## Bus Factor: %d\n\n", stats.BusFactor))
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

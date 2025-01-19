package output

import (
	"fmt"
	"strings"

	"github.com/erdemkosk/gitness/internal/models"
)

type ConsoleFormatter struct{}

func (f *ConsoleFormatter) Format(stats *models.RepositoryStats) (string, error) {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("\nRepo: %s/%s\n", stats.Owner, stats.Repo))
	output.WriteString(fmt.Sprintf("Bus Factor: %d\n\n", stats.BusFactor))
	output.WriteString("Contributors:\n")
	output.WriteString("------------------\n")

	for _, c := range stats.Contributors {
		output.WriteString(fmt.Sprintf("%s: %d commits (%.1f%%)\n",
			c.Name, c.Commits, c.Percentage))
	}

	return output.String(), nil
}

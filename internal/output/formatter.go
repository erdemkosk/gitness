package output

import (
	"fmt"

	"github.com/erdemkosk/gitness/internal/models"
)

type Formatter interface {
	Format(*models.RepositoryStats) (string, error)
}

type FormatterFactory struct {
	formatters map[string]Formatter
}

func NewFormatterFactory() *FormatterFactory {
	factory := &FormatterFactory{
		formatters: make(map[string]Formatter),
	}

	// Register default formatters
	factory.Register("console", &ConsoleFormatter{})
	factory.Register("json", &JSONFormatter{})
	factory.Register("markdown", &MarkdownFormatter{})

	return factory
}

func (f *FormatterFactory) Register(name string, formatter Formatter) {
	f.formatters[name] = formatter
}

func (f *FormatterFactory) GetFormatter(format string) (Formatter, bool) {
	fmt.Println(format)
	formatter, exists := f.formatters[format]
	return formatter, exists
}

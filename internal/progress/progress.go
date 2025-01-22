package progress

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func NewProgressBar(description string) *ProgressBar {
	bar := progressbar.NewOptions64(
		-1, // indeterminate mode
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50), // Sabit genişlik
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s...", description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionSetWriter(os.Stderr),
	)

	return &ProgressBar{
		bar: bar,
	}
}

func (p *ProgressBar) SetTotal(total int64) {
	p.bar.ChangeMax64(total)
}

func (p *ProgressBar) Increment() {
	p.bar.Add64(1)
}

func (p *ProgressBar) Finish() {
	p.bar.Finish()
	p.bar.Clear() // İşlem bittiğinde progress bar'ı temizle
}

func (p *ProgressBar) Clear() {
	p.bar.Clear()
}

package progress

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressBarOptions struct {
	Description string
	Width       int
	SpinnerType int
	ThrottleMS  int
	Theme       *ProgressBarTheme
}

type ProgressBarTheme struct {
	Saucer        string
	SaucerHead    string
	SaucerPadding string
	BarStart      string
	BarEnd        string
}

var DefaultTheme = &ProgressBarTheme{
	Saucer:        "[green]=[reset]",
	SaucerHead:    "[green]>[reset]",
	SaucerPadding: " ",
	BarStart:      "[",
	BarEnd:        "]",
}

var DefaultOptions = &ProgressBarOptions{
	Width:       50,
	SpinnerType: 14,
	ThrottleMS:  65,
	Theme:       DefaultTheme,
}

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func NewProgressBar(description string) *ProgressBar {
	return NewProgressBarWithOptions(&ProgressBarOptions{
		Description: description,
		Width:       DefaultOptions.Width,
		SpinnerType: DefaultOptions.SpinnerType,
		ThrottleMS:  DefaultOptions.ThrottleMS,
		Theme:       DefaultOptions.Theme,
	})
}

func NewProgressBarWithOptions(opts *ProgressBarOptions) *ProgressBar {
	if opts.Theme == nil {
		opts.Theme = DefaultTheme
	}

	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(opts.Width),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s...", opts.Description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        opts.Theme.Saucer,
			SaucerHead:    opts.Theme.SaucerHead,
			SaucerPadding: opts.Theme.SaucerPadding,
			BarStart:      opts.Theme.BarStart,
			BarEnd:        opts.Theme.BarEnd,
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(opts.SpinnerType),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(time.Duration(opts.ThrottleMS)*time.Millisecond),
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
}

func (p *ProgressBar) Clear() {
	p.bar.Clear()
}

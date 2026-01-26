package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
)

func NewSpinner(_ string) *spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &s
}

func NewProgress(_ int) *progress.Model {
	p := progress.New(progress.WithDefaultGradient())
	return &p
}

package ui

import "github.com/charmbracelet/lipgloss"

var fileStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Width(100) //nolint:gomnd,mnd

var removedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#ffebeb")).
	Width(100) //nolint:gomnd,mnd

var addedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#f4f8f1")).
	Width(100) //nolint:gomnd,mnd

var hr = lipgloss.NewStyle().
	Background(lipgloss.Color("#f6f6f6")).
	Width(100) //nolint:gomnd,mnd

var titleStyle = lipgloss.NewStyle().Padding(0, 0)

var infoStyle = titleStyle.Copy()

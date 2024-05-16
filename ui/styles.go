package ui

import "github.com/charmbracelet/lipgloss"

var fileStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Width(100)

var removedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#ffe8e7")).
	Width(100)

var addedStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#f0f8ec")).
	Width(100)

var hr = lipgloss.NewStyle().
	Background(lipgloss.Color("#f6f6f6")).
	Width(100)

var titleStyle = lipgloss.NewStyle().Padding(0, 0)

var infoStyle = func() lipgloss.Style {
	return titleStyle.Copy()
}()

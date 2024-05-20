package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jshawl/abcd/parser"
)

type File struct {
	parser.File
}

type Block struct {
	OldRange string
	NewRange string
	Lines    []string
}

func NewFile(file parser.File) File {
	return File{
		file,
	}
}

var lineNumberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaa"))

func (m File) View(viewportWidth int) string {
	var content strings.Builder
	content.WriteString(fileStyle.Width(viewportWidth).Render(m.Name))
	content.WriteString("\n")
	for blockI, block := range m.Blocks {
		for _, line := range block.Lines {
			largestLineNumber := fmt.Sprintf("%d", block.LargestLineNumber)
			fmtString := fmt.Sprintf("%%%dd ", lipgloss.Width(largestLineNumber))
			lineNumber := lineNumberStyle.Render(fmt.Sprintf(fmtString, line.Number))
			width := viewportWidth - lipgloss.Width(largestLineNumber) - 1
			content.WriteString(lineNumber)
			if strings.HasPrefix(line.Content, "-") {
				content.WriteString(removedStyle.Width(width).Render(line.Content))
			} else if strings.HasPrefix(line.Content, "+") {
				content.WriteString(addedStyle.Width(width).Render(line.Content))
			} else {
				content.WriteString(line.Content)
			}
			content.WriteString("\n")
		}
		if blockI < len(m.Blocks)-1 {
			content.WriteString(hr.Width(viewportWidth).Render("···"))
			content.WriteString("\n")
		}
	}
	return content.String()
}

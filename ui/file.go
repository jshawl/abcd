package ui

import (
	"fmt"
	"strconv"
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
		newStart := strings.Split(block.NewRange, ",")
		newLineCounter, _ := strconv.Atoi(newStart[0])
		oldStart := strings.Split(block.OldRange, ",")
		oldLineCounter, _ := strconv.Atoi(oldStart[0])
		for i, line := range block.Lines {

			if strings.HasPrefix(line, "-") {
				lineNumber := lineNumberStyle.Render(fmt.Sprintf("%d ", i+oldLineCounter))
				width := viewportWidth - lipgloss.Width(lineNumber)
				content.WriteString(lineNumber)
				newLineCounter--
				content.WriteString(removedStyle.Width(width).Render(line))
			} else if strings.HasPrefix(line, "+") {
				lineNumber := lineNumberStyle.Render(fmt.Sprintf("%d ", i+newLineCounter))
				width := viewportWidth - lipgloss.Width(lineNumber)
				content.WriteString(lineNumber)
				oldLineCounter--
				content.WriteString(addedStyle.Width(width).Render(line))
			} else {
				lineNumber := lineNumberStyle.Render(fmt.Sprintf("%d ", i+newLineCounter))
				content.WriteString(lineNumber)
				content.WriteString(line)
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

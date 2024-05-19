package ui

import (
	"strings"

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

func (m File) View(viewportWidth int) string {
	var content strings.Builder
	content.WriteString(fileStyle.Width(viewportWidth).Render(m.Name))
	content.WriteString("\n")
	for blockI, block := range m.Blocks {
		for _, line := range block.Lines {
			if strings.HasPrefix(line, "-") {
				content.WriteString(removedStyle.Width(viewportWidth).Render(line))
			} else if strings.HasPrefix(line, "+") {
				content.WriteString(addedStyle.Width(viewportWidth).Render(line))
			} else {
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

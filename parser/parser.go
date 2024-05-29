package parser

import (
	"bufio"
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Diff struct {
	Files []File
}

type File struct {
	Name   string
	Blocks []Block
}

type Block struct {
	LargestLineNumber int
	Lines             []Line
	NewStart          int
	NewEnd            int
	OldStart          int
	OldEnd            int
}

type Line struct {
	Number  int
	Content string
}

func ParseLine(line string) (string, error) {
	r := regexp.MustCompile(`^(diff|index|---|\+\+\+|@@)`)
	if !r.MatchString(line) {
		return line, nil
	}

	return "", MatchError()
}

func ParseBlock(line string) (Block, error) {
	r := regexp.MustCompile(`^@@ -([0-9]+),([0-9]+) \+([0-9]+),?([0-9]*) @@`)
	matches := r.FindAllStringSubmatch(line, -1)

	if len(matches) != 1 {
		return Block{}, MatchError()
	}

	OldStart, _ := strconv.Atoi(matches[0][1])
	OldLines, _ := strconv.Atoi(matches[0][2])
	OldEnd := OldStart + OldLines
	NewStart, _ := strconv.Atoi(matches[0][3])
	NewLines, _ := strconv.Atoi(matches[0][4])
	NewEnd := NewStart + NewLines

	return Block{
		OldStart:          OldStart,
		OldEnd:            OldEnd,
		NewStart:          NewStart,
		NewEnd:            NewEnd,
		LargestLineNumber: 0,
		Lines:             []Line{},
	}, nil
}

func ParseFile(line string) (File, error) {
	r := regexp.MustCompile(`^(?:\-\-\- a\/|\+\+\+ b\/)(.*)`)
	matches := r.FindAllStringSubmatch(line, -1)

	if len(matches) != 1 {
		return File{}, MatchError()
	}

	return File{Name: matches[0][1], Blocks: []Block{}}, nil
}

func ParseDiff(lines string) (Diff, error) {
	var parsedLines []string

	diff := Diff{
		Files: []File{},
	}
	sc := bufio.NewScanner(strings.NewReader(lines))

	for sc.Scan() {
		parsedLines = append(parsedLines, sc.Text())
	}

	var (
		oldLineCounter int
		newLineCounter int
	)

	for _, value := range parsedLines {
		file, _ := ParseFile(value)

		if file.Name == "" && len(diff.Files) == 0 {
			continue
		}

		if file.Name != "" {
			fileExists := slices.ContainsFunc(diff.Files, func(f File) bool {
				return f.Name == file.Name
			})
			// the file name might have been identified with
			// `--- a/(.*)` or `+++ b/(.*)` but if identified with both,
			// don't append a second file with the same name.
			if !fileExists {
				diff.Files = append(diff.Files, file)
			}
		}

		lastFile := &diff.Files[len(diff.Files)-1]
		block, err := ParseBlock(value)

		if err == nil {
			newLineCounter = block.NewStart
			oldLineCounter = block.OldStart
			block.LargestLineNumber = max(block.NewEnd, block.OldEnd)
			lastFile.Blocks = append(lastFile.Blocks, block)
		}

		blocks := lastFile.Blocks
		line, err := ParseLine(value)

		if err == nil {
			structuredLine := Line{Content: line, Number: 0}

			if strings.HasPrefix(line, "-") {
				structuredLine.Number = oldLineCounter
				newLineCounter--
			} else if strings.HasPrefix(line, "+") {
				structuredLine.Number = newLineCounter
				oldLineCounter--
			} else {
				structuredLine.Number = newLineCounter
			}

			blocks[len(blocks)-1].Lines = append(blocks[len(blocks)-1].Lines, structuredLine)
			newLineCounter++
			oldLineCounter++
		}
	}

	return diff, nil
}

var errMatch = errors.New("match not found")

func MatchError() error {
	return errMatch
}

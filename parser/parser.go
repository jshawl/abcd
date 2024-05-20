package parser

import (
	"bufio"
	"errors"
	"regexp"
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

func parseLine(line string) (string, error) {
	r := regexp.MustCompile(`^(diff|index|---|\+\+\+|@@)`)
	if !r.MatchString(line) {
		return line, nil
	}
	return "", errors.New("no match")
}

func parseBlock(line string) (Block, error) {
	r := regexp.MustCompile(`^@@ -([0-9]+),([0-9]+) \+([0-9]+),?([0-9]*) @@`)
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) != 1 {
		return Block{}, errors.New("match not found")
	}

	OldStart, _ := strconv.Atoi(matches[0][1])
	OldLines, _ := strconv.Atoi(matches[0][2])
	OldEnd := OldStart + OldLines
	NewStart, _ := strconv.Atoi(matches[0][3])
	NewLines, _ := strconv.Atoi(matches[0][4])
	NewEnd := NewStart + NewLines

	return Block{OldStart: OldStart, OldEnd: OldEnd, NewStart: NewStart, NewEnd: NewEnd}, nil
}

func parseFile(line string) (File, error) {
	r := regexp.MustCompile(`^\+\+\+ b\/(.*)`)
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) != 1 {
		return File{}, errors.New("match not found")
	}
	return File{Name: matches[0][1]}, nil
}

func ParseDiff(lines string) (Diff, error) {
	diff := Diff{}
	var parsedLines []string
	sc := bufio.NewScanner(strings.NewReader(lines))
	for sc.Scan() {
		parsedLines = append(parsedLines, sc.Text())
	}
	var (
		oldLineCounter int
		newLineCounter int
	)
	for _, v := range parsedLines {
		file, _ := parseFile(v)
		if file.Name == "" && len(diff.Files) == 0 {
			continue
		}
		if file.Name != "" {
			diff.Files = append(diff.Files, file)
		}
		lastFile := &diff.Files[len(diff.Files)-1]
		block, err := parseBlock(v)
		if err == nil {
			newLineCounter = block.NewStart
			oldLineCounter = block.OldStart
			block.LargestLineNumber = block.NewEnd
			lastFile.Blocks = append(lastFile.Blocks, block)
		}
		blocks := lastFile.Blocks
		line, err := parseLine(v)
		if err == nil {

			l := Line{Content: line}

			if strings.HasPrefix(line, "-") {
				l.Number = oldLineCounter
				newLineCounter--
			} else if strings.HasPrefix(line, "+") {
				l.Number = newLineCounter
				oldLineCounter--
			} else {
				l.Number = newLineCounter
			}
			blocks[len(blocks)-1].Lines = append(blocks[len(blocks)-1].Lines, l)
			newLineCounter++
			oldLineCounter++
		}
	}
	return diff, nil
}

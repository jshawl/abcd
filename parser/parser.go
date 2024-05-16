package parser

import (
	"bufio"
	"errors"
	"regexp"
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
	OldRange string
	NewRange string
	Lines    []string
}

func parseLine(line string) (string, error) {
	r := regexp.MustCompile(`^(diff|index|---|\+\+\+|@@)`)
	if !r.MatchString(line) {
		return line, nil
	}
	return "", errors.New("no match")
}

func parseBlock(line string) (Block, error) {
	r := regexp.MustCompile(`@@ -([0-9]+,[0-9]+) \+([0-9]+,?[0-9]*) @@`)
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) != 1 {
		return Block{}, errors.New("match not found")
	}
	oldRange := matches[0][1]
	newRange := matches[0][2]
	return Block{OldRange: oldRange, NewRange: newRange}, nil
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
	for _, v := range parsedLines {
		file, _ := parseFile(v)
		if file.Name == "" && len(diff.Files) == 0 {
			continue
		}
		if file.Name != "" {
			diff.Files = append(diff.Files, file)
		}
		lastFile := &diff.Files[len(diff.Files)-1]
		block, _ := parseBlock(v)
		if block.OldRange != "" {
			lastFile.Blocks = append(lastFile.Blocks, block)
		}
		blocks := lastFile.Blocks
		line, err := parseLine(v)
		if err == nil {
			blocks[len(blocks)-1].Lines = append(blocks[len(blocks)-1].Lines, line)
		}
	}
	return diff, nil
}

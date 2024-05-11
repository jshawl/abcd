package main

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
	r := regexp.MustCompile(`diff --git a/([\w\.]+) b/`)
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) != 1 {
		return File{}, errors.New("match not found")
	}
	return File{Name: matches[0][1]}, nil
}

func parseDiff(lines string) (Diff, error) {
	diff := Diff{}
	var parsedLines []string
	sc := bufio.NewScanner(strings.NewReader(lines))
	for sc.Scan() {
		parsedLines = append(parsedLines, sc.Text())
	}
	for _, v := range parsedLines {
		file, _ := parseFile(v)
		if file.Name != "" {
			diff.Files = append(diff.Files, file)
		}
	}
	return diff, nil
}

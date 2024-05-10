package main

import (
	"errors"
	"regexp"
)

type Diff struct {
	Files []string
}

func parseFile(line string) (string, error) {
	r := regexp.MustCompile(`diff --git a/([\w\.]+) b/`)
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) != 1 {
		return "", errors.New("match not found")
	}
	return matches[0][1], nil
}

// func parse(diff string) Diff {
// 	re := regexp.MustCompile("^diff")
// 	output := re.Split(string(diff), -1)
// 	return Diff{Files: []string{""}}
// }

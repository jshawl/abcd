package main

import (
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {
	expected := "file"
	actual, _ := parseFile("diff --git a/file b/file")
	if expected != actual {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual))
	}
}

func TestParseFileWithExtension(t *testing.T) {
	expected := "file.txt"
	actual, _ := parseFile("diff --git a/file.txt b/file.txt")
	if expected != actual {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual))
	}
}

func TestParseFileOnNonFile(t *testing.T) {
	_, err := parseFile("not a line with a filename")
	if err == nil {
		t.Fatalf("Expected error, got match")
	}
}

// func TestParseOneFileOneBlock(t *testing.T) {
// 	diff, _ := os.ReadFile("./test/one-file-one-block.diff")
// 	parsed := parse(string(diff))

// 	if parsed.Files[0] != "file" {
// 		t.Fatalf("Expected actual")
// 	}
// }

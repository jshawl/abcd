package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	expected := "file"
	actual, _ := parseFile("diff --git a/file b/file")
	if expected != actual.Name {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual))
	}
}

func TestParseFileWithExtension(t *testing.T) {
	expected := "file.txt"
	actual, _ := parseFile("diff --git a/file.txt b/file.txt")
	if expected != actual.Name {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual))
	}
}

func TestParseFileOnNonFile(t *testing.T) {
	_, err := parseFile("not a line with a filename")
	if err == nil {
		t.Fatalf("Expected error, got match")
	}
}

func TestParseBlockWithAddedLines(t *testing.T) {
	actual, _ := parseBlock("@@ -0,0 +1 @@")
	if actual.OldRange != "0,0" {
		t.Fatalf("Expected old range to be '0,0'")
	}
	if actual.NewRange != "1" {
		t.Fatalf("Expected new range to be '1'")
	}
}

func TestParseBlockWithChangedLines(t *testing.T) {
	actual, _ := parseBlock("@@ -1,4 +1,4 @@")
	if actual.OldRange != "1,4" {
		t.Fatalf("Expected old range to be '1,4'")
	}
	if actual.NewRange != "1,4" {
		t.Fatalf("Expected new range to be '1,4'")
	}
}

func TestParseDiff(t *testing.T) {
	contents, _ := os.ReadFile("./test/one-file-one-block.diff")
	actual, _ := parseDiff(string(contents))
	if len(actual.Files) != 1 {
		t.Fatalf("Expected 1 file")
	}
	if actual.Files[0].Name != "file" {
		t.Fatalf("Expected 1 file name to be 'file'")
	}
}

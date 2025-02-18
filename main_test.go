package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputfile  = "./testdata/test1.md"
	goldenfile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {

	input, err := os.ReadFile(inputfile)

	if err != nil {
		t.Fatal(err)
	}
	result, err := parseContent(input, "")

	expected, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}

}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer
	if err := run(inputfile, "", &mockStdOut, true); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(resultFile)

	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenfile)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}

	os.Remove(resultFile)
}

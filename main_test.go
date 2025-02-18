package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputfile  = "./testdata/test1.md"
	resultfile = "test1.md.html"
	goldenfile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {

	input, err := os.ReadFile(inputfile)

	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)

	expected, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected,result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}

}

func TestRun(t *testing.T) {
	if err := run(inputfile); err != nil {
		t.Fatal(err)
	}

	result, err := os.ReadFile(resultfile)

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

	os.Remove(resultfile)
}

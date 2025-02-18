package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>`

	footer = `</body>
</html>`
)

func main() {

	file := flag.String("file", "", "File to preview")

	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filename string) error {
	input, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	parse := parseContent(input)

	temp, err := os.CreateTemp("", "mdp*.html")

	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}
	outputName := temp.Name()
	return saveToHtml(outputName, parse)
}

func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)

	html := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(html)
	buffer.WriteString(footer)
	return buffer.Bytes()
}

func saveToHtml(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

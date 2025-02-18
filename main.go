package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*file, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filename string, out io.Writer, skipPreview bool) error {
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
	fmt.Fprintln(out, outputName)
	if err := saveToHtml(outputName, parse); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outputName)
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

func preview(fname string) error {
	cName := ""
	cParams := []string{}
	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// Append filename to parameters slice
	cParams = append(cParams, fname)
	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	// Open the file using default program
	return exec.Command(cPath, cParams...).Run()
}

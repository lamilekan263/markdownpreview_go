package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {

	file := flag.String("file", "", "File to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*file, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, tFname)

	temp, err := os.CreateTemp("", "mdp*.html")

	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}
	outputName := temp.Name()
	fmt.Fprintln(out, outputName)
	if err := saveToHtml(outputName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	defer os.Remove(outputName)
	return preview(outputName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)

	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)

	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer

	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
	err = exec.Command(cPath, cParams...).Run()

	time.Sleep(2 * time.Second)
	return err
}

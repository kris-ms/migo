package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type MDConverter struct {
	inputPath  string
	outputPath string
    templatePath string
}

// Takes an absolute .md path and attempts to output rendered html to the [MDConverter].outputPath.
// Returns the m.outputPath if successful.
func (m *MDConverter) Convert() (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
            extension.GFM,
            highlighting.NewHighlighting(
                highlighting.WithStyle("vs"),
            ),
        ),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

    source, err := m.getSource()
    if err != nil {
        return "", fmt.Errorf("get source from m.inputpath, %v: %v", source, err)
    }

	var convertBuffer bytes.Buffer
	err = md.Convert(source, &convertBuffer)
	if err != nil {
		return "", fmt.Errorf("convert md %v: %v", source, md)
	}
	htmlOutput := make([]byte, convertBuffer.Len())
    br, err := convertBuffer.Read(htmlOutput)
	if err != nil {
		return "", fmt.Errorf("read converted htmlOutput to buf %v: %v", htmlOutput, err)
	}
	fmt.Println(br, "bytes read from htmlOutput buffer.")

    out, err := m.makeOutput(htmlOutput)

    if err != nil {
        return "", fmt.Errorf("failed make output %v: %v", out, err)
    }

	return m.outputPath, nil
}

func (m *MDConverter) getSource() ([]byte, error) {
	if !path.IsAbs(m.inputPath) {
		return make([]byte, 0), errors.New("Invalid inputPath, please provide absolute path to .md file.")
	}

	inputFile, err := os.OpenFile(m.inputPath, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return make([]byte, 0), fmt.Errorf("open file %v: %v", inputFile, err)
	}

    defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return make([]byte, 0), fmt.Errorf("file info %v: %v", fileInfo, err)
	}

	source := make([]byte, fileInfo.Size())
	br, err := inputFile.Read(source)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("read %v to source buffer: %v, bytes read: %d", source, err, br)
	}
	fmt.Println(br, "bytes read from source.")

    return source, nil
}

func (m *MDConverter) makeOutput(convertedHtml []byte) ([]byte, error) {
	output, err := m.makeOutputDir()
	if err != nil {
		return make([]byte, 0), fmt.Errorf("create file %v: %v, %v", m.outputPath, err, output)
	}
    defer output.Close()

    tmpl := TemplateMaker{
        inputPath: m.templatePath,
    }

    type TestData struct {
        PageTitle string
        ConvertedHtml string
    }

    testdata := TestData{
        PageTitle: "Testing template from md",
        ConvertedHtml: string(convertedHtml),
    }

    out, err := tmpl.Execute(testdata)

    if err != nil {
        return make([]byte, 0), fmt.Errorf("fail tmpl execute %v, %v: %v", testdata, out, err)
    }

	bw, err := output.Write(out)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("write out to file %v, %v: %v", out, m.outputPath, err)
	}
	fmt.Println(bw, "bytes written to", m.outputPath)

    return out, nil
}

func (m *MDConverter) makeOutputDir() (*os.File, error) {
    if err := os.MkdirAll(filepath.Dir(m.outputPath), 0770); err != nil {
        return nil, err
    }
    return os.Create(m.outputPath)
}

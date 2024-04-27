package main

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"text/template"

	"github.com/yosssi/gohtml"
)

type TemplateMaker struct {
    inputPath string
}

func (t *TemplateMaker) Execute(data any) ([]byte, error) {
    var rdbuf bytes.Buffer
    err := t.getTemplate().Execute(&rdbuf, data)
    if err != nil {
        return make([]byte, 0), fmt.Errorf("executing tmpl: %v", err)
    }
    out := make([]byte, rdbuf.Len())
    bytesread, err := rdbuf.Read(out)
    if err != nil {
        return make([]byte, 0), fmt.Errorf("reading tmpl: %v", err)
    }
    fmt.Printf("read %d bytes from executed template\n", bytesread)
    return gohtml.FormatBytes(out), nil
}

func (t *TemplateMaker) getTemplate() *template.Template {
    if !path.IsAbs(t.inputPath) {
       panic(errors.New("invalid path, provide absolute path to .html template file")) 
    }
    tmpl, err := template.ParseFiles(t.inputPath)
    if err != nil {
        panic(err)
    }
    return tmpl
}

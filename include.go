package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Include struct {
	includeDir string
}

func (i *Include) Copy(buildDir string) error {
	builddir, err := filepath.Abs(buildDir)
	if err != nil {
		return fmt.Errorf("bad builddir %v: %v", builddir, err)
	}
	incpath, err := filepath.Abs(i.includeDir)
	if err != nil {
		return fmt.Errorf("bad include path %v: %v", incpath, err)
	}

	includeFiles := make([]string, 0)
	outputFiles := make([]string, 0)

	collect := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error collecting include file %v, %v", d, err)
		}
		if d.IsDir() == false {
			includeFiles = append(includeFiles, path)
			outputFiles = append(outputFiles, buildDir+strings.TrimPrefix(path, i.includeDir))
		}
		return nil
	}

	if err := filepath.WalkDir(incpath, collect); err != nil {
		return fmt.Errorf("failed walkdir %v: %v", incpath, err)
	}

	for i, file := range includeFiles {
		input, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("error open input %v: %v", file, err)
		}
		if err := os.MkdirAll(filepath.Dir(outputFiles[i]), 0777); err != nil {
			return fmt.Errorf("error mkdir output %v: %v", outputFiles[i], err)
		}
		output, err := os.Create(outputFiles[i])
		if err != nil {
			return fmt.Errorf("error create output file %v: %v", output, err)
		}
		if bw, err := io.Copy(output, input); err != nil {
			return fmt.Errorf("error copying files %v to %v, %v bytes written: %v", input, output, bw, err)
		}
	}
	return nil
}

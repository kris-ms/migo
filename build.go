package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Builder struct {
	workDir string
}

func (b *Builder) Build() error {
	var normalizedWorkDir string
	if !strings.HasSuffix(b.workDir, "/") {
		normalizedWorkDir = b.workDir + "/"
	} else {
		normalizedWorkDir = b.workDir
	}
	var outdir string
	outdir = normalizedWorkDir + "build/"

	if err := os.RemoveAll(outdir); err != nil {
		return fmt.Errorf("failed to clean build: %v", err)
	}

	files := []string{}
	err := filepath.WalkDir(b.workDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir failed: %v", err)
	}
	mdFiles := []string{}
	templateFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file, ".md") {
			mdFiles = append(mdFiles, file)
		}
		if strings.HasSuffix(file, ".html") {
			templateFiles = append(templateFiles, file)
		}
	}

	for _, fileToConvert := range mdFiles {

		filename := filepath.Base(fileToConvert)
		fileNoExt := strings.TrimSuffix(filename, ".md")

		matchIdx := slices.IndexFunc(templateFiles, func(ele string) bool {
			tmplFileName := filepath.Base(ele)
			tmplFileNoExt := strings.TrimSuffix(tmplFileName, ".html")
			return tmplFileNoExt == fileNoExt
		})

		if matchIdx == -1 {
			return fmt.Errorf("no template match for %v: %v", fileNoExt, templateFiles)
		}

		outpath := outdir + fileNoExt + ".html"

		md := MDConverter{
			inputPath:    fileToConvert,
			outputPath:   outpath,
			templatePath: templateFiles[matchIdx],
		}

		out, err := md.Convert()
		if err != nil {
			return fmt.Errorf("bad convert to md %v, %v: %v", out, md, err)
		}
	}
	include := Include{
		includeDir: normalizedWorkDir + "include/",
	}

	if err := include.Copy(outdir); err != nil {
		return fmt.Errorf("failed to copy %v to %v: %v", include, outdir, err)
	}
	return nil
}

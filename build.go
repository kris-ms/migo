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
		return err
	}

	files := []string{}
	err := filepath.WalkDir(b.workDir, func(path string, d fs.DirEntry, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
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
	fmt.Printf("md files: %v\n", mdFiles)
	fmt.Printf("template files: %v\n", templateFiles)

	for _, fileToConvert := range mdFiles {

		filename := filepath.Base(fileToConvert)
		fileNoExt := strings.TrimSuffix(filename, ".md")

		fmt.Println(fileNoExt)

		matchIdx := slices.IndexFunc(templateFiles, func(ele string) bool {
			tmplFileName := filepath.Base(ele)
			tmplFileNoExt := strings.TrimSuffix(tmplFileName, ".html")
			return tmplFileNoExt == fileNoExt
		})

		if matchIdx == -1 {
			return fmt.Errorf("match templates %v idx %d: %v, outdir: %v", templateFiles, matchIdx, fileNoExt, outdir)
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
		fmt.Println(out)
	}
	include := Include{
		includeDir: normalizedWorkDir + "include/",
	}

	if err := include.Copy(outdir); err != nil {
		return err
	}
	return nil
}

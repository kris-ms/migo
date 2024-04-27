package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "migo",
		Usage: "build a static website from markdown files and html templates",
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "start building inside of provided directory",
				Action: func(cCtx *cli.Context) error {
					absWorkdir, err := filepath.Abs(cCtx.Args().First())
					if err != nil {
						return fmt.Errorf("invalid path error %v: %v", absWorkdir, err)
					}

					builder := Builder{
						workDir: absWorkdir,
					}

					fmt.Printf("Attempting to build in directory %v\n", builder.workDir)
					start := time.Now().UnixMilli()

					if err := builder.Build(); err != nil {
						return fmt.Errorf("failed to build with error: %v", err)
					}

					end := time.Now().UnixMilli() - start

					fmt.Printf(
						"Successfully build to directory %v, time elapsed: %v\n",
						builder.workDir+"/build",
						fmt.Sprintf("%vms", end),
					)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

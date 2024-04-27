package main

import (
    "fmt"
    "os"
    "log"
    "path/filepath"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name: "migo",
        Usage: "build a static website from markdown files and html templates",
        Commands: []*cli.Command{
            {
                Name: "build",
                Aliases: []string{"b"},
                Usage: "start building inside of provided directory",
                Action: func(cCtx *cli.Context) error {
                    absWorkdir, err := filepath.Abs(cCtx.Args().First())
                    if err != nil {
                        return err
                    }
                    
                    builder := Builder{
                        workDir: absWorkdir,
                    }

                    fmt.Printf("Attempting to build in directory %v\n", builder.workDir)

                    err = builder.Build()
                    return err
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

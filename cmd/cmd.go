package cmd

import (
	"encoding/json"
	"fmt"
	relgen "github.com/bajankristof/relgen/internal"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

var (
	Version     = "0.0.1"
	Description = "RelGen is a command-line tool that helps you automate version bumps and changelog creation for any of your projects by using Conventional Commits and your Git history."
)

var (
	ConfigFlag        = "config"
	PreReleaseFlag    = "pre-release"
	BuildMetadataFlag = "build-metadata"
	VersionPrefixFlag = "version-prefix"
	DryRunFlag        = "dry-run"
)

func Start() error {
	app := &cli.App{
		Name:        "RelGen",
		HelpName:    "relgen",
		Version:     Version,
		Description: Description,
		Commands:    []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      ConfigFlag,
				Usage:     "path to the configuration file to use (command line flags will take precedence)",
				Value:     "./relgenrc.json",
				Aliases:   []string{"c"},
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:  PreReleaseFlag,
				Usage: "generate a pre-release version with the specified tag",
				Value: "",
			},
			&cli.StringFlag{
				Name:  BuildMetadataFlag,
				Usage: "set the build-metadata of the generated version",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  VersionPrefixFlag,
				Usage: "generate the release version with a 'v' prefix (e.g.: v1.0.0)",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  DryRunFlag,
				Usage: "print the generated release to the standard output",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			cfg, err := relgen.ReadConfig(path.Join(cwd, ctx.String(ConfigFlag)))
			if err != nil {
				return err
			}

			fmt.Println(cfg)

			if ctx.IsSet(PreReleaseFlag) {
				cfg.PreRelease = ctx.String(PreReleaseFlag)
			}

			if ctx.IsSet(BuildMetadataFlag) {
				cfg.BuildMetadata = ctx.String(BuildMetadataFlag)
			}

			if ctx.IsSet(VersionPrefixFlag) {
				cfg.VersionPrefix = ctx.Bool(VersionPrefixFlag)
			}

			repo, err := git.PlainOpen(cwd)
			if err != nil {
				return err
			}

			builder := relgen.NewReleaseBuilder(repo, cfg)
			rel, err := builder.Build()
			if err != nil {
				return err
			}

			output, _ := json.Marshal(rel)
			fmt.Println(string(output))

			return nil
		},
	}

	return app.Run(os.Args)
}

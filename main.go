package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/insomnimus/gerp/cmd"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed completion_help/bash.txt
	bashCompletionMsg string

	//go:embed completion_help/powershell.txt
	psCompletionMsg string

	//go:embed completion_help/zsh.txt
	zshCompletionMsg string

	//go:embed complete/*
	completions embed.FS
)

const VERSION = "0.3.0"

const helpMsg = `gerp, match regular expressions
usage:
	gerp [options] <pattern> [file...]
options are:
	-i, --ignore-case: do a case insensitive search
	-v, --invert: print lines not matching the pattern
	-m, --match: only display matches
	-n, --no-header: do not print any header info
	-q, --quiet: do not print errors
	-d, --hidden: do not ignore hidden files and directories
	-g, --glob=<pattern>: use gerps globbing engine for searching files (not required on windows)
	--generate-completion: generate shell autocompletions for bash, powershell or zsh
	--help-completions: display helpabout installing shell completions
	-V, --version: show the gerp version installed
`

func generateCompletions(sh string) error {
	data, err := completions.ReadFile("complete/" + strings.ToLower(sh))
	if err != nil {
		return fmt.Errorf("unrecognized shell for autocomplete. known shells are [bash, powershell, zsh]")
	}
	fmt.Println(string(data))
	return nil
}

func bf(short, long, desc string, target *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        long,
		Usage:       desc,
		Destination: target,
		Aliases:     []string{short},
	}
}

func showCompletionHelp(sh string) {
	switch strings.ToLower(sh) {
	case "bash":
		fmt.Println(bashCompletionMsg)
	case "powershell":
		fmt.Println(psCompletionMsg)
	case "zsh":
		fmt.Println(zshCompletionMsg)
	default:
		log.Fatalf("%s: unrecognized shell. available shells are bash, powershell and zsh", sh)
	}
}

func showUsage(_ *cli.Context, err error, _ bool) error {
	log.Fatalf("%s\nuse with --help for the usage", err)
	return nil
}

func showVersion() {
	fmt.Printf("gerp version %s\n", VERSION)
	os.Exit(0)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("")
	if len(os.Args) <= 1 {
		log.Println("gerp, match regular expressions\nuse with --help for the usage")
		os.Exit(0)
	}

	var (
		helpCompletions string
		flagComplete    string
		flagVersion     bool
	)

	opt := new(cmd.Cmd)
	run := func(c *cli.Context) error {
		if helpCompletions != "" {
			showCompletionHelp(helpCompletions)
			return nil
		}
		if flagVersion {
			showVersion()
			return nil
		}
		if flagComplete != "" {
			return generateCompletions(flagComplete)
		}

		opt.Pattern = c.Args().First()
		opt.Args = c.Args().Tail()
		if err := opt.Process(); err != nil {
			return err
		}
		return opt.Run()
	}

	app := &cli.App{
		OnUsageError:          showUsage,
		CustomAppHelpTemplate: helpMsg,
		Name:                  "gerp",
		Version:               VERSION,
		ArgsUsage:             "gerp [OPTIONS] [FILE...]",
		HideHelpCommand:       true,
		HideVersion:           true,
		Usage:                 "match regular expressions",
		Action:                run,
		Flags: []cli.Flag{
			bf("i", "ignore-case", "ignore case while matching", &opt.IgnoreCase),
			bf("v", "invert", "print lines not matching the pattern", &opt.Invert),
			bf("m", "match", "only print text matching the pattern (not the whole line)", &opt.Match),
			bf("d", "hidden", "do not ignore files and directories starting with '.'", &opt.Hidden),
			bf("q", "quiet", "do not report non-fatal errors", &opt.Quiet),
			bf("n", "no-header", "do not print headers", &opt.NoHeader),
			&cli.StringFlag{
				Name:        "help-completions",
				Usage:       "print help about shell autocompletions (bash, powershell or zsh)",
				Destination: &helpCompletions,
			},
			&cli.StringFlag{
				Name:        "generate-completion",
				Usage:       "generate a shell completion script for one of [bash, powershell, zsh]",
				Destination: &flagComplete,
			},
			&cli.StringFlag{
				Name:        "glob",
				Aliases:     []string{"g"},
				Usage:       "use gerp's globbing engine to search files (not required on windows",
				Destination: &opt.Glob,
			},
			bf("V", "version", "show gerp version and exit", &flagVersion),
		},
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

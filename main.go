package main

import (
	"fmt"
	"gerp/cmd"
	"log"
	"os"
	"runtime"
)

func showAbout() {
	fmt.Println("gerp, match regular expressions\nrun with --help for the usage")
	os.Exit(0)
}

func showHelp() {
	fmt.Println(`gerp, match regular expressions
usage:
	gerp [options] <pattern> [file...]
options are:
	-i, --ignore-case: do a case insensitive search
	-v, --invert: print lines not matching the pattern
	-m, --match: only display matches
	-n, --no-header: do not print any header info
	-q, --quiet: do not print errors
	-d, --hidden: do not ignore hidden files and directories
	--: indicate that the rest of the arguments are file names
	-h, --help: show this message and exit
	`)
	if runtime.GOOS == "windows" {
		fmt.Println("file names can be glob patterns supporting double star")
	} else {
		fmt.Println("\t--glob: specify a glob pattern to use")
	}
	os.Exit(0)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("")
	if len(os.Args) <= 1 {
		showAbout()
	}
	c, err := cmd.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	if c.FlagH {
		showHelp()
	}
	// fmt.Println(c.Pattern)
	// check if stdin is piped
	if fi, err := os.Stdin.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			c.RunStdin()
			return
		}
	}

	err = c.Run()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

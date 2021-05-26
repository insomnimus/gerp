package cmd

import (
	"regexp"
)

type Cmd struct {
	re               *regexp.Regexp
	Glob             string
	Pattern          string
	Args             []string
	Files            []string
	Quiet            bool
	NoHeader         bool
	Match            bool
	Hidden           bool
	Invert           bool
	isReg            bool
	filesAreFiltered bool
	IgnoreCase       bool
}

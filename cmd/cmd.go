package cmd

import (
	"regexp"
)

var (
	reFlag = regexp.MustCompile(`^\-[mvidqn]+$`)
)

type Cmd struct {
	IgnoreCase, Invert bool
	Hidden, Quiet      bool
	NoHeader, Match    bool

	Pattern string
	Args    []string
	Files   []string
	Glob    string

	isReg            bool
	filesAreFiltered bool
	re               *regexp.Regexp
}

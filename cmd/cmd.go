package cmd

import (
	"regexp"
	//"fmt"
)

var (
	reFlag = regexp.MustCompile(`^\-[mvidqn]+$`)
)

type Cmd struct {
	FlagI, FlagD, FlagM, FlagV, FlagQ, FlagN bool
	FlagVersion, FlagH                       bool
	Files                                    []string
	Pattern                                  string

	isReg            bool
	filesAreFiltered bool
	re               *regexp.Regexp
}

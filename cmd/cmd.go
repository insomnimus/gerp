package cmd

import "regexp"

var (
	reFlag = regexp.MustCompile(`^\-[mvidn]+$`)
)

type Cmd struct {
	FlagH, FlagI, FlagD, FlagM, FlagV, FlagQ, FlagN bool
	Files                                           []string
	Pattern                                         string

	isReg            bool
	filesAreFiltered bool
	re               *regexp.Regexp
}

//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"strings"

	"github.com/mattn/go-zglob"
)

func (cmd *Cmd) Process() error {
	if cmd.Pattern == "" {
		return fmt.Errorf("missing required argument: pattern")
	}
	if cmd.Glob != "" {
		cmd.Args = append(cmd.Args, cmd.Glob)
	}
	tmp := make(map[string]struct{})
	for _, f := range cmd.Args {
		if strings.ContainsAny(f, "*{?") {
			fs, err := zglob.Glob(f)
			if err != nil {
				return err
			}
			for _, x := range fs {
				tmp[x] = struct{}{}
			}
			continue
		}
		tmp[f] = struct{}{}
	}
	cmd.Files = make([]string, 0, len(tmp))
	for f := range tmp {
		cmd.Files = append(cmd.Files, f)
	}
	return nil
}

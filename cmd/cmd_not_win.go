//go:build !windows
// +build !windows

package cmd

import (
	"fmt"
	"github.com/mattn/go-zglob"
)

func (cmd *Cmd) Process() error {
	if cmd.Pattern == "" {
		return fmt.Errorf("missing required argument: pattern")
	}
	cmd.Files = cmd.Args
	if cmd.Glob != "" {
		files, err := zglob.Glob(cmd.Glob)
		if err != nil {
			return err
		}
		cmd.Files = append(cmd.Files, files...)
	}
	return nil
}

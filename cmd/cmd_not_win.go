//go:build !windows
// +build !windows

package cmd

import (
	"fmt"
	"github.com/mattn/go-zglob"
	"strings"
)

func Parse(args []string) (*Cmd, error) {
	if len(args) == 0 {
		return &Cmd{FlagH: true}, nil
	}

	set := make(map[rune]struct{})
	var a string
	cmd := Cmd{
		Files:            make([]string, 0, len(args)-1),
		filesAreFiltered: true,
	}
	setflag := func(c rune, b bool, s string) {
		if _, ok := set[c]; !ok {
			set[c] = struct{}{}
			switch c {
			case 'n':
				cmd.FlagN = true
			case 'd':
				cmd.FlagD = b
			case 'i':
				cmd.FlagI = b
			case 'm':
				cmd.FlagM = b
			case 'q':
				cmd.FlagQ = b
			case 'v':
				cmd.FlagV = b
			}
		} else {
			cmd.Files = append(cmd.Files, s)
		}
	}

LOOP:
	for i := 0; i < len(args); i++ {
		a = args[i]
		if a == "" {
			continue
		}
		if a[0] == '-' {
			if reFlag.MatchString(a) {
				for _, c := range a[1:] {
					switch c {
					case 'q':
						cmd.FlagQ = true
					case 'd':
						cmd.FlagD = true
					case 'i':
						cmd.FlagI = true
					case 'm':
						cmd.FlagM = true
					case 'v':
						cmd.FlagV = true
					}
					set[c] = struct{}{}
				}
				continue
			}
			switch a {
			case "--no-header", "--no-header=true":
				setflag('n', true, a)
			case "--no-header=false":
				setflag('n', false, a)
			case "--quiet", "--quiet=true":
				setflag('q', true, a)
			case "--quiet=false":
				setflag('q', false, a)
			case "-h", "--help":
				return &Cmd{FlagH: true}, nil
			case "--":
				cmd.Files = append(cmd.Files, args[i+1:]...)
				break LOOP
			case "--hidden", "--hidden=true":
				setflag('d', true, a)
			case "--hidden=false":
				setflag('d', false, a)
			case "--ignore-case", "--ignore-case=true":
				setflag('i', true, a)
			case "--ignore-=false":
				setflag('i', false, a)
			case "--invert", "--invert=true":
				setflag('r', true, a)
			case "--invert=false":
				setflag('r', false, a)
			case "--match", "--match=true":
				setflag('m', true, a)
			case "--match=false":
				setflag('m', false, a)
			case "--glob":
				if i+1 >= len(args) {
					return nil, fmt.Errorf("the --glob flag is set but the value is missing")
				}
				i++
				fs, err := zglob.Glob(args[i])
				if err != nil {
					return nil, err
				}
				cmd.filesareFiltered = false
				cmd.Files = append(cmd.Files, fs...)
			case "--pattern":
				if i+1 >= len(args) {
					return nil, fmt.Errorf("%s flag is set but the value is missing", a)
				}
				if cmd.Pattern != "" {
					cmd.Files = append(cmd.Files, cmd.Pattern)
				}
				i++
				cmd.Pattern = args[i]
			default:
				if strings.HasPrefix(a, "--glob=") {
					fs, err := zglob.Glob(strings.TrimPrefix(a, "--glob="))
					if err != nil {
						return nil, err
					}
					cmd.filesAreFiltered = false
					cmd.Files = append(cmd.Files, fs...)
					continue
				}
				if strings.HasPrefix(a, "--pattern=") {
					if cmd.Pattern != "" {
						cmd.Files = append(cmd.Files, cmd.Pattern)
					}
					cmd.Pattern = strings.TrimPrefix(a, "--pattern")
					continue
				}

				if cmd.Pattern == "" {
					cmd.Pattern = a
				} else {
					cmd.Files = append(cmd.Files, a)
				}
			}
			continue
		}
		if cmd.Pattern == "" {
			cmd.Pattern = a
		} else {
			cmd.Files = append(cmd.Files, a)
		}
	}
	if cmd.Pattern == "" {
		return nil, fmt.Errorf("the pattern is missing")
	}
	return &cmd, nil
}

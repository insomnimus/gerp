//go:build windows

package cmd

import (
	"fmt"
	"strings"

	"github.com/mattn/go-zglob"
)

func Parse(args []string) (*Cmd, error) {
	if len(args) == 0 {
		return &Cmd{FlagH: true}, nil
	}

	set := make(map[rune]struct{})
	var a string
	cmd := Cmd{
		Files: make([]string, 0, len(args)-1),
	}

	setflag := func(c rune, b bool, s string) {
		if _, ok := set[c]; !ok {
			set[c] = struct{}{}
			switch c {
			case 'n':
				cmd.FlagN = b
			case 'q':
				cmd.FlagQ = b
			case 'd':
				cmd.FlagD = b
			case 'i':
				cmd.FlagI = b
			case 'm':
				cmd.FlagM = b
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
					case 'n':
						cmd.FlagN = true
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
			case "--version":
				return &Cmd{FlagVersion: true}, nil
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
				if strings.HasPrefix(a, "--pattern=") {
					if cmd.Pattern != "" {
						cmd.Files = append(cmd.Files, cmd.Pattern)
					}
					cmd.Pattern = strings.TrimPrefix(a, "--pattern")
					continue
				}

				return nil, fmt.Errorf("unknown command line option %q", a)
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
	tmp := make(map[string]struct{})
	for _, f := range cmd.Files {
		if strings.ContainsAny(f, "*{?") {
			fs, err := zglob.Glob(f)
			if err != nil {
				return nil, err
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

	return &cmd, nil
}

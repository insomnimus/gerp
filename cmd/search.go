package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var reHidden = regexp.MustCompile(`\a\.|[\\/]\.`)

func (c *Cmd) Run() error {
	if c.IgnoreCase ||
		regexp.MustCompile(`[^\w\d_\s-]`).MatchString(c.Pattern) {
		c.isReg = true
		var err error
		if c.IgnoreCase {
			c.re, err = regexp.Compile("(?i)" + c.Pattern)
		} else {
			c.re, err = regexp.Compile(c.Pattern)
		}
		if err != nil {
			return err
		}
	}
	if fi, err := os.Stdin.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			c.NoHeader = true
			return c.RunStdin()
		}
	}
	if len(c.Files) == 0 {
		var err error
		c.Files, err = collectFiles(".", c.Hidden)
		if err != nil {
			return err
		}
		if len(c.Files) == 0 {
			return nil
		}
		c.filesAreFiltered = true
	}

	workerN := 8
	if n := runtime.NumCPU(); n > workerN {
		workerN = n
	}
	if workerN > len(c.Files) {
		workerN = len(c.Files)
	}

	jobN := len(c.Files)
	jobs := make(chan string, jobN)
	results := make(chan struct{}, jobN)

	for i := 0; i < workerN; i++ {
		go c.worker(jobs, results)
	}
	for _, f := range c.Files {
		jobs <- f
	}
	close(jobs)
	for i := 0; i < jobN; i++ {
		<-results
	}
	return nil
}

func (c *Cmd) worker(jobs <-chan string, results chan<- struct{}) {
	for j := range jobs {
		c.search(j)
		results <- struct{}{}
	}
}

func (c *Cmd) search(name string) {
	// do not read exe, object, bin files
	if !c.filesAreFiltered && (strings.HasSuffix(name, ".exe") ||
		strings.HasSuffix(name, ".bin") ||
		strings.HasSuffix(name, ".o")) {
		return
	}
	// do not read if file is hidden
	if !c.filesAreFiltered && !c.Hidden && reHidden.MatchString(name) {
		return
	}

	f, err := os.Open(name)
	if err != nil {
		if !c.Quiet {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				log.Println(err)
			}
		}
		return
	}
	defer f.Close()

	var (
		scanner = bufio.NewScanner(f)
		i       uint32
		s       string
		ok      bool
		buf     strings.Builder
	)

	for scanner.Scan() {
		i++
		if scanner.Err() != nil {
			return
		}
		s, ok = c.eval(scanner.Text())
		if !ok {
			continue
		}
		if c.NoHeader {
			fmt.Fprintln(&buf, s)
		} else {
			fmt.Fprintf(&buf, "%-6d|  %s\n", i, s)
		}
	}

	if buf.Len() > 0 {
		if !c.NoHeader {
			fmt.Printf("# %s\n%s", name, buf.String())
		} else {
			fmt.Print(buf.String())
		}
	}
}

func (c *Cmd) eval(s string) (str string, ok bool) {
	if !c.isReg {
		if strings.Contains(s, c.Pattern) {
			if c.Invert {
				if !c.Match {
					return
				}
				return strings.ReplaceAll(s, c.Pattern, ""), true
			}
			// is a match, flagV false
			if c.Match {
				return c.Pattern, true
			}
			return s, true
		}
		// does not contain
		if c.Invert {
			return s, true
		}
		return
	}
	// is regex
	// if not flagM, just check if regex matches
	if !c.Match {
		if c.re.MatchString(s) {
			if c.Invert {
				return
			}
			return s, true
		}
		// not a match, FlagM false
		if c.Invert {
			return s, true
		}
		return
	}

	// FlagM true

	if !c.Invert {
		text := c.re.FindString(s)
		return text, text != ""
	}

	text := c.re.ReplaceAllString(s, "")
	return text, text != ""
}

func (c *Cmd) RunStdin() error {
	if c.IgnoreCase ||
		regexp.MustCompile(`[^\w\d_\s-]`).MatchString(c.Pattern) {
		c.isReg = true
		var err error
		if c.IgnoreCase {
			c.re, err = regexp.Compile("(?i)" + c.Pattern)
		} else {
			c.re, err = regexp.Compile(c.Pattern)
		}
		if err != nil {
			return err
		}
	}

	var (
		scanner = bufio.NewScanner(os.Stdin)
		i       uint32
		s       string
		ok      bool
	)

	for scanner.Scan() {
		i++
		if scanner.Err() != nil {
			return nil
		}
		s, ok = c.eval(scanner.Text())
		if !ok {
			continue
		}
		if c.NoHeader {
			fmt.Println(s)
		} else {
			fmt.Printf("%-6d:    %s\n", i, s)
		}
	}
	return nil
}

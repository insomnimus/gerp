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
	if c.FlagI ||
		regexp.MustCompile(`[^\w\d_\s-]`).MatchString(c.Pattern) {
		c.isReg = true
		var err error
		c.re, err = regexp.Compile(c.Pattern)
		if err != nil {
			return err
		}
	}
	if len(c.Files) == 0 {
		var err error
		c.Files, err = collectFiles("./", c.FlagD)
		if err != nil {
			return err
		}
		if len(c.Files) == 0 {
			return nil
		}
		c.filesAreFiltered = true
	}

	workerN := 4
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
	// do not read if file is hidden
	if !c.filesAreFiltered && !c.FlagD && reHidden.MatchString(name) {
		return
	}

	f, err := os.Open(name)
	if err != nil {
		if !c.FlagQ {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				log.Println(err)
			}
		}
		return
	}
	defer f.Close()

	var (
		scanner        = bufio.NewScanner(f)
		i       uint32 = 1
		s       string
		ok      bool
		isFirst = true
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
		if !c.FlagN {
			if isFirst {
				isFirst = false
				fmt.Printf("--%s--\n", name)
			}
			if !c.FlagV {
				fmt.Printf("%-4d:    %s\n", i, s)
			} else {
				fmt.Println(s)
			}
		} else {
			fmt.Println(s)
		}
	}
}

func (c *Cmd) eval(s string) (str string, ok bool) {
	if !c.isReg {
		if strings.Contains(s, c.Pattern) {
			if c.FlagV {
				if !c.FlagM {
					return
				}
				return strings.ReplaceAll(s, c.Pattern, ""), true
			}
			// is a match, flagV false
			if c.FlagM {
				return c.Pattern, true
			}
			return s, true
		}
		// does not contain
		if c.FlagV {
			return s, true
		}
		return
	}
	// is regex
	// if not flagM, just check if regex matches
	if !c.FlagM {
		if c.re.MatchString(s) {
			if c.FlagV {
				return
			}
			return s, true
		}
		// not a match, FlagM false
		if c.FlagV {
			return s, true
		}
		return
	}

	// FlagM true

	if !c.FlagV {
		text := c.re.FindString(s)
		return text, text != ""
	}

	text := c.re.ReplaceAllString(s, "")
	return text, text != ""
}

func (c *Cmd) RunStdin() {
	var (
		scanner        = bufio.NewScanner(os.Stdin)
		i       uint32 = 1
		s       string
		ok      bool
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
		if !c.FlagN {

			if !c.FlagV {
				fmt.Printf("%-4d:    %s\n", i, s)
			} else {
				fmt.Println(s)
			}
		} else {
			fmt.Println(s)
		}
	}
}

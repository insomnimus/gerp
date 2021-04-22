package cmd

import (
	"github.com/mattn/go-zglob/fastwalk"
	"io/fs"
	"os"
	"sync"
)

func collectFiles(dir string, includeHidden bool) ([]string, error) {
	var mux sync.Mutex
	var files []string
	err := fastwalk.FastWalk(dir, func(path string, mode os.FileMode) error {
		if !includeHidden && reHidden.MatchString(path) {
			if mode == os.ModeDir {
				return fs.SkipDir
			}
			return nil
		}
		if mode != os.ModeDir {
			mux.Lock()
			files = append(files, path)
			mux.Unlock()
		}
		return nil
	})
	return files, err
}

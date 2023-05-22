package fsutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/must"
	"github.com/jopbrown/gobase/strutil"
)

func AppDir() string {
	exePath := must.Value(os.Executable())
	return filepath.Dir(exePath)
}

func WorkDirWithMatcher(root string, matcher strutil.Matcher, fn fs.WalkDirFunc) error {
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.ErrorAt(err)
		}

		rel, err := filepath.Rel(root, p)
		if err != nil {
			return errors.ErrorAt(err)
		}

		if !matcher.MatchString(filepath.ToSlash(rel)) {
			return nil
		}

		err = fn(p, d, err)
		if err != nil {
			return errors.ErrorAt(err)
		}

		return nil
	})
	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

func ListWithMatcher(root string, matcher strutil.Matcher) ([]string, error) {
	paths := make([]string, 0, 10)
	err := WorkDirWithMatcher(root, matcher, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.ErrorAt(err)
		}

		paths = append(paths, path)
		return nil
	})

	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return paths, nil
}

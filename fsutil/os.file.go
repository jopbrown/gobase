package fsutil

import (
	"os"
	"path/filepath"
	"time"

	"github.com/jopbrown/gobase/errors"
)

func OpenFileWrite(elem ...string) (*os.File, error) {
	fname := filepath.Join(elem...)
	err := os.MkdirAll(filepath.Dir(fname), 0755)
	if err != nil {
		return nil, errors.ErrorAt(err, "unable to create folder of file: ", fname)
	}

	f, err := os.Create(fname)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return f, nil
}

func OpenFileRead(elem ...string) (*os.File, error) {
	fname := filepath.Join(elem...)
	f, err := os.Open(fname)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return f, nil
}

func OpenFileAppend(elem ...string) (*os.File, error) {
	fname := filepath.Join(elem...)
	err := os.MkdirAll(filepath.Dir(fname), 0755)
	if err != nil {
		return nil, errors.ErrorAt(err, "unable to create folder of file: ", fname)
	}

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return f, nil
}

func FileTouch(elem ...string) error {
	fname := filepath.Join(elem...)
	err := os.MkdirAll(filepath.Dir(fname), 0755)
	if err != nil {
		return errors.ErrorAt(err, "unable to create folder of file: ", fname)
	}

	if ExistsFile(fname) {
		now := time.Now()
		err := os.Chtimes(fname, now, now)
		if err != nil {
			return errors.ErrorAt(err, "unable to change modify time of file: ", fname)
		}
	} else {
		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return errors.ErrorAt(err, "unable to open file: ", fname)
		}
		f.Close()
	}

	return nil
}

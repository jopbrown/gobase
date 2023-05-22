package fsutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jopbrown/gobase/errors"
)

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func ExistsDir(name string) bool {
	finfo, err := os.Stat(name)
	if err != nil {
		return false
	}

	return finfo.IsDir()
}

func ExistsFile(name string) bool {
	finfo, err := os.Stat(name)
	if err != nil {
		return false
	}

	return finfo.Mode().IsRegular()
}

// Copy file or folder from src to dst.
// the behavior same to shell's cp command.
func Copy(dst, src string) error {
	var err error
	if ExistsDir(src) {
		err = CopyDir(dst, src)
	} else if ExistsFile(src) {
		err = CopyFile(dst, src)
	} else {
		err = errors.Error("source not exist: ", src)
	}

	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

func CopyFile(dst, src string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.ErrorAt(err, "src not exist: ", src)
	}

	dst = filepath.ToSlash(dst)
	if ExistsDir(dst) || dst[len(dst)-1] == '/' {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return errors.ErrorAt(err, "unable to create folder of dst: ", dst)
	}

	fin, err := os.Open(src)
	if err != nil {
		return errors.ErrorAt(err, "unable to open src: ", src)
	}

	defer fin.Close()

	fout, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return errors.ErrorAt(err, "unable to open dst: ", dst)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return errors.ErrorAt(err, "unable to copy file stream")
	}

	return nil
}

func CopyDir(dstDir, srcDir string) error {
	err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if srcDir == path {
			return nil
		}
		src := path
		relsrc, errrel := filepath.Rel(srcDir, path)
		if errrel != nil {
			return nil
		}
		dst := filepath.Join(dstDir, relsrc)
		if ExistsDir(src) {
			CopyDir(dst, src)
		} else {
			CopyFile(dst, src)
		}
		return nil
	})

	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

// Move file or folder from src to dst.
// The behavior same to shell's mv command.
// function signature same to func Move(srcs..., dst) error.
func Move(args ...string) error {
	argCnt := len(args)
	if argCnt < 2 {
		return errors.Errorf("move: wrong number of args: %v", args)
	}

	srcList, dst := args[:argCnt-1], args[argCnt-1]
	dst = filepath.ToSlash(dst)

	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return errors.ErrorAt(err, "unable to create folder of dst: ", dst)
	}

	manySrc := false
	isDstDir := ExistsDir(dst) || dst[len(dst)-1] == '/'

	for _, src := range srcList {
		if !Exists(src) {
			if err != nil {
				return errors.ErrorAt(err, "src not exist: ", src)
			}
		}

		if manySrc && !isDstDir {
			return errors.Errorf("too many src move to single dst: %s %s", src, dst)
		}

		manySrc = true

		if ExistsDir(src) && dst[len(dst)-1] == '/' {
			dstDir := dst
			if dstDir[len(dstDir)-1] == '/' {
				dstDir = filepath.Join(dst, filepath.Base(src))
			}
			err = os.Rename(src, dstDir)
			if err != nil {
				return errors.ErrorAtf(err, "unable to move %s to %s", src, dstDir)
			}
			continue
		}

		dstFile := dst
		if isDstDir {
			dstFile = filepath.Join(dst, filepath.Base(src))
		}

		err = os.Rename(src, dstFile)
		if err != nil {
			return errors.ErrorAtf(err, "unable to move %s to %s", src, dstFile)
		}
	}

	return nil
}

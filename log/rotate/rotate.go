package rotate

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/djherbis/times"
	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/jopbrown/gobase/must"
)

type Writer struct {
	maxSize int64

	interval       time.Duration
	lastRotateTime time.Time

	fd *os.File

	fpath       string
	perm        fs.FileMode
	written     int64
	rotateCount int
}

func OpenFile(name string, interval time.Duration, maxSize int64) (*Writer, error) {
	f, err := fsutil.OpenFileAppend(name)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}
	return NewWriter(f, interval, maxSize), nil
}

func NewWriter(fd *os.File, interval time.Duration, maxSize int64) *Writer {
	w := &Writer{}
	w.interval = interval
	w.maxSize = maxSize
	w.fd = fd
	finfo := must.Value(fd.Stat())
	w.written = finfo.Size()
	w.perm = finfo.Mode().Perm()
	w.fpath = fd.Name()

	w.lastRotateTime = getCreateTime(fd)
	return w
}

func getCreateTime(fd *os.File) time.Time {
	stat := must.Value(times.StatFile(fd))
	if stat.HasBirthTime() {
		return stat.BirthTime()
	}
	createTime := stat.AccessTime()
	if stat.ModTime().Sub(createTime) < 0 {
		createTime = stat.ModTime()
		if stat.HasChangeTime() && stat.ChangeTime().Sub(createTime) < 0 {
			createTime = stat.ChangeTime()
		}
	}
	return createTime
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.fd.Write(p)
	if err != nil {
		return
	}

	w.written += int64(n)

	err = w.doRotate()
	return
}

func (w *Writer) Close() error {
	return w.fd.Close()
}

func (w *Writer) doRotate() error {
	var now time.Time
	needRotate := false
	if w.interval > 0 {
		now = time.Now()
		if now.Sub(w.lastRotateTime) > w.interval {
			needRotate = true
		}
	}

	if w.maxSize > 0 && w.written > w.maxSize {
		now = time.Now()
		needRotate = true
	}

	if needRotate {
		err := w.rotateFile(now)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) rotateFile(now time.Time) error {
	w.rotateCount++
	noext, ext := filePathSplitByExt(w.fpath)
	// backupPath := noext +  + ext
	backupPath := fmt.Sprintf("%s.%s_%02d%s", noext, now.Format("20060102_150405"), w.rotateCount, ext)

	err := w.fd.Close()
	if err != nil {
		errors.ErrorAt(err)
	}

	err = os.Rename(w.fpath, backupPath)
	if err != nil {
		errors.ErrorAt(err)
	}

	w.fd, err = os.OpenFile(w.fpath, os.O_CREATE|os.O_WRONLY, w.perm)
	if err != nil {
		errors.ErrorAt(err)
	}
	w.lastRotateTime = now

	return nil
}

func filePathSplitByExt(fpath string) (noext, ext string) {
	for i := len(fpath) - 1; i >= 0 && !os.IsPathSeparator(fpath[i]); i-- {
		if fpath[i] == '.' {
			return fpath[:i], fpath[i:]
		}
	}
	return fpath, ""
}

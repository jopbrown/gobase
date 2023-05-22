package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jopbrown/gobase/errors"
	"golang.org/x/exp/slog"
)

type Logger struct {
	mu        sync.Mutex
	out       io.Writer
	format    LoggerFormat
	buf       []byte
	isDiscard atomic.Bool

	prefix   string
	minLevel Level
	maxLevel Level
	verbose  int
}

func NewLogger(out io.Writer, minLevel, maxLevel Level) *Logger {
	l := &Logger{}
	l.out = out
	l.minLevel = minLevel
	l.maxLevel = maxLevel
	l.format = DefaultLoggerFormat()
	if out == io.Discard {
		l.isDiscard.Store(true)
	}
	return l
}

func NewLoggerWithFormat(out io.Writer, minLevel, maxLevel Level, format LoggerFormat) *Logger {
	l := NewLogger(out, minLevel, maxLevel)
	l.format = format
	return l
}

func (l *Logger) Format() LoggerFormat {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.format
}

func (l *Logger) SetFormat(format LoggerFormat) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = format
}

func (l *Logger) GetWriter(level Level) io.Writer {
	if !l.enabled(level) {
		return io.Discard
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
	l.isDiscard.Store(w == io.Discard)
}

func (l *Logger) Clone() *Logger {
	newl := NewLogger(l.out, l.minLevel, l.maxLevel)
	newl.verbose = l.verbose
	newl.format = l.format
	newl.prefix = l.prefix
	return newl
}

func (l *Logger) V(v int) ILogger {
	newl := l.Clone()
	newl.verbose = l.verbose + v
	return newl
}

func (l *Logger) With(prefix string) ILogger {
	newl := l.Clone()
	newl.prefix = path.Join(l.prefix, prefix)
	return newl
}

func (l *Logger) S(json bool) *slog.Logger {
	h := newSLoggerHandler(l, json)
	s := slog.New(h)

	return s
}

func (l *Logger) enabled(level Level) bool {
	if l.isDiscard.Load() {
		return false
	}

	if l.verbose > int(globalVerbose.Load()) {
		return false
	}

	if level == LevelAll {
		return true
	}

	if level == LevelNone {
		return false
	}

	if level < l.minLevel {
		return false
	}

	if level > l.maxLevel {
		return false
	}

	return true
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func getCaller(skip int) *runtime.Frame {
	callers := make([]uintptr, 1)
	n := runtime.Callers(skip, callers[:])
	if n < 1 {
		return nil
	}
	frame, _ := runtime.CallersFrames(callers).Next()
	return &frame
}

func (l *Logger) output(calldepth int, level Level, msg string) error {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = l.buf[:0]

	if l.format.AddLevel {
		lvlStr := level.String()
		l.buf = append(l.buf, lvlStr...)
		l.buf = append(l.buf, bytes.Repeat([]byte{' '}, 5-len(lvlStr)+1)...)
	}

	if l.format.AddVerbose {
		l.buf = append(l.buf, 'V')
		itoa(&l.buf, l.verbose, 1)
		l.buf = append(l.buf, ' ')
	}

	if l.format.AddDateTime {
		if l.format.DateTimeFormat.AddDate {
			year, month, day := now.Date()
			itoa(&l.buf, year, 4)
			l.buf = append(l.buf, '/')
			itoa(&l.buf, int(month), 2)
			l.buf = append(l.buf, '/')
			itoa(&l.buf, day, 2)
			l.buf = append(l.buf, ' ')
		}
		if l.format.DateTimeFormat.AddTime {
			hour, min, sec := now.Clock()
			itoa(&l.buf, hour, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, min, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, sec, 2)
			if l.format.DateTimeFormat.AddMicroseconds {
				l.buf = append(l.buf, '.')
				itoa(&l.buf, now.Nanosecond()/1e3, 6)
			}
			l.buf = append(l.buf, ' ')
		}
	}

	if l.format.AddSource || l.format.AddCaller {
		frame := getCaller(calldepth + 1)

		if l.format.AddSource {
			file := frame.File
			fileDepthCount := 0
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					newFile := file[i+1:]
					fileDepthCount++
					if l.format.SourceDepth > 0 && fileDepthCount >= l.format.SourceDepth {
						file = newFile
						break
					}
				}
			}

			l.buf = append(l.buf, file...)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, frame.Line, -1)
			l.buf = append(l.buf, ": "...)
		}

		if l.format.AddCaller {
			if l.format.LongCaller {
				l.buf = append(l.buf, frame.Function...)
			} else {
				l.buf = append(l.buf, path.Base(frame.Function)...)
			}
			l.buf = append(l.buf, ' ')
		}
	}

	if l.format.AddPrefix && len(l.prefix) > 0 {
		l.buf = append(l.buf, l.prefix...)
		l.buf = append(l.buf, ' ')
	}

	l.buf = append(l.buf, msg...)
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	if err != nil {
		return errors.ErrorAt(err)
	}
	return nil
}

func (l *Logger) Print(a ...any) {
	if !l.enabled(LevelAll) {
		return
	}
	fmt.Fprint(l.out, a...)
}

func (l *Logger) Printf(format string, a ...any) {
	if !l.enabled(LevelAll) {
		return
	}
	fmt.Fprintf(l.out, format, a...)
}

func (l *Logger) Println(a ...any) {
	if !l.enabled(LevelAll) {
		return
	}
	fmt.Fprintln(l.out, a...)
}

func (l *Logger) Printlnf(format string, a ...any) {
	if !l.enabled(LevelAll) {
		return
	}
	fmt.Fprintf(l.out, format, a...)
	io.WriteString(l.out, "\n")
}

func (l *Logger) Debug(a ...any) {
	if !l.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprint(a...)
	l.output(3, LevelDebug, msg)
}

func (l *Logger) Debugf(format string, a ...any) {
	if !l.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	l.output(3, LevelDebug, msg)
}

func (l *Logger) Info(a ...any) {
	if !l.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprint(a...)
	l.output(3, LevelInfo, msg)
}

func (l *Logger) Infof(format string, a ...any) {
	if !l.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	l.output(3, LevelInfo, msg)
}

func (l *Logger) Warn(a ...any) {
	if !l.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprint(a...)
	l.output(3, LevelWarn, msg)
}

func (l *Logger) Warnf(format string, a ...any) {
	if !l.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	l.output(3, LevelWarn, msg)
}

func (l *Logger) Error(a ...any) {
	if !l.enabled(LevelError) {
		return
	}
	msg := fmt.Sprint(a...)
	l.output(3, LevelError, msg)
}

func (l *Logger) Errorf(format string, a ...any) {
	if !l.enabled(LevelError) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	l.output(3, LevelError, msg)
}

func (l *Logger) ErrorAt(err error, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprint(a...))
	if !l.enabled(LevelError) {
		return err
	}
	l.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func (l *Logger) ErrorAtf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprintf(format, a...))
	if !l.enabled(LevelError) {
		return err
	}
	l.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func (l *Logger) Fatal(a ...any) {
	if l.enabled(LevelFatal) {
		msg := fmt.Sprint(a...)
		l.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, a ...any) {
	if l.enabled(LevelFatal) {
		msg := fmt.Sprintf(format, a...)
		l.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func (l *Logger) Panic(a ...any) {
	msg := fmt.Sprint(a...)
	if l.enabled(LevelPanic) {
		l.output(3, LevelPanic, msg)
	}
	panic(msg)
}

func (l *Logger) Panicf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if l.enabled(LevelPanic) {
		l.output(3, LevelPanic, msg)
	}
	panic(msg)
}

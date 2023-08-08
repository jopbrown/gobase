package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"log/slog"

	"github.com/jopbrown/gobase/errors"
)

type ILogger interface {
	enabled(level Level) bool
	output(calldepth int, level Level, msg string) error

	GetWriter(level Level) io.Writer
	V(v int) ILogger
	With(prefix string) ILogger
	S(json bool) *slog.Logger

	Print(a ...any)
	Printf(format string, a ...any)
	Println(a ...any)
	Printlnf(format string, a ...any)
	Debug(a ...any)
	Debugf(format string, a ...any)
	Info(a ...any)
	Infof(format string, a ...any)
	Warn(a ...any)
	Warnf(format string, a ...any)
	Error(a ...any)
	Errorf(format string, a ...any)
	ErrorAt(err error, a ...any) error
	ErrorAtf(err error, format string, a ...any) error
	Fatal(a ...any)
	Fatalf(format string, a ...any)
	Panic(a ...any)
	Panicf(format string, a ...any)
}

var (
	globalVerbose atomic.Int32
	globalLogger  ILogger = DefaultLogger(false)
)

func init() {
	globalVerbose.Store(0)
}

func DefaultLogger(debug bool) ILogger {
	minLevel := LevelInfo
	if debug {
		minLevel = LevelDebug
	}
	return NewLogger(os.Stderr, minLevel, LevelFatal)
}

func ConsoleLogger(debug bool) ILogger {
	minLevel := LevelInfo
	if debug {
		minLevel = LevelDebug
	}
	otherLog := NewLoggerWithFormat(os.Stdout, minLevel, LevelInfo, SimpleLoggerFormat())
	errLog := NewLoggerWithFormat(os.Stderr, LevelWarn, LevelFatal, FileLoggerFormat())
	return NewTeeLogger(errLog, otherLog)
}

func FileLogger(w io.Writer, format LoggerFormat, debug bool) ILogger {
	minLevel := LevelInfo
	if debug {
		minLevel = LevelDebug
	}
	return NewLoggerWithFormat(w, minLevel, LevelFatal, format)
}

func SetGlobalVerbose(v int) int {
	old := int(globalVerbose.Swap(int32(v)))
	return old
}

func SetGlobalLogger(l ILogger) {
	globalLogger = l
}

func GetWriter(level Level) io.Writer {
	return globalLogger.GetWriter(level)
}

func V(v int) ILogger {
	return globalLogger.V(v)
}

func With(prefix string) ILogger {
	return globalLogger.With(prefix)
}

func S(json bool) *slog.Logger {
	return globalLogger.S(json)
}

func Print(a ...any) {
	globalLogger.Print(a...)
}

func Printf(format string, a ...any) {
	globalLogger.Printf(format, a...)
}

func Println(a ...any) {
	globalLogger.Println(a...)
}

func Printlnf(format string, a ...any) {
	globalLogger.Printlnf(format, a...)
}

func Debug(a ...any) {
	if !globalLogger.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprint(a...)
	globalLogger.output(3, LevelDebug, msg)
}

func Debugf(format string, a ...any) {
	if !globalLogger.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	globalLogger.output(3, LevelDebug, msg)
}

func Info(a ...any) {
	if !globalLogger.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprint(a...)
	globalLogger.output(3, LevelInfo, msg)
}

func Infof(format string, a ...any) {
	if !globalLogger.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	globalLogger.output(3, LevelInfo, msg)
}

func Warn(a ...any) {
	if !globalLogger.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprint(a...)
	globalLogger.output(3, LevelWarn, msg)
}

func Warnf(format string, a ...any) {
	if !globalLogger.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	globalLogger.output(3, LevelWarn, msg)
}

func Error(a ...any) {
	if !globalLogger.enabled(LevelError) {
		return
	}
	msg := fmt.Sprint(a...)
	globalLogger.output(3, LevelError, msg)
}

func Errorf(format string, a ...any) {
	if !globalLogger.enabled(LevelError) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	globalLogger.output(3, LevelError, msg)
}

func ErrorAt(err error, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprint(a...))
	if !globalLogger.enabled(LevelError) {
		return err
	}
	globalLogger.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func ErrorAtf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprintf(format, a...))
	if !globalLogger.enabled(LevelError) {
		return err
	}
	globalLogger.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func Fatal(a ...any) {
	if globalLogger.enabled(LevelFatal) {
		msg := fmt.Sprint(a...)
		globalLogger.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	if globalLogger.enabled(LevelFatal) {
		msg := fmt.Sprintf(format, a...)
		globalLogger.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func Panic(a ...any) {
	msg := fmt.Sprint(a...)
	if globalLogger.enabled(LevelPanic) {
		globalLogger.output(3, LevelPanic, msg)
	}
	panic(msg)
}

func Panicf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if globalLogger.enabled(LevelPanic) {
		globalLogger.output(3, LevelPanic, msg)
	}
	panic(msg)
}

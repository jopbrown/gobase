package log

import (
	"fmt"
	"io"
	"os"

	"log/slog"

	"github.com/jopbrown/gobase/errors"
)

type TeeLogger struct {
	loggers []ILogger
}

func NewTeeLogger(loggers ...ILogger) *TeeLogger {
	tee := &TeeLogger{}
	tee.loggers = loggers
	return tee
}

func (tee *TeeLogger) GetWriter(level Level) io.Writer {
	ws := make([]io.Writer, 0, 1)
	for _, l := range tee.loggers {
		w := l.GetWriter(level)
		if w != io.Discard {
			ws = append(ws, w)
		}
	}

	if len(ws) == 0 {
		return io.Discard
	}

	return io.MultiWriter(ws...)
}

func (tee *TeeLogger) V(v int) ILogger {
	loggers := make([]ILogger, 0, len(tee.loggers))
	for _, l := range tee.loggers {
		loggers = append(loggers, l.V(v))
	}

	return NewTeeLogger(loggers...)
}

func (tee *TeeLogger) With(prefix string) ILogger {
	loggers := make([]ILogger, 0, len(tee.loggers))
	for _, l := range tee.loggers {
		loggers = append(loggers, l.With(prefix))
	}

	return NewTeeLogger(loggers...)
}

func (tee *TeeLogger) S(json bool) *slog.Logger {
	h := newSTeeLoggerHandler(tee, json)
	s := slog.New(h)

	return s
}

func (tee *TeeLogger) enabled(level Level) bool {
	for _, l := range tee.loggers {
		if l.enabled(level) {
			return true
		}
	}
	return false
}

func (tee *TeeLogger) Print(a ...any) {
	if !tee.enabled(LevelAll) {
		return
	}
	fmt.Fprint(tee.GetWriter(LevelAll), a...)
}

func (tee *TeeLogger) Printf(format string, a ...any) {
	if !tee.enabled(LevelAll) {
		return
	}
	fmt.Fprintf(tee.GetWriter(LevelAll), format, a...)
}

func (tee *TeeLogger) Println(a ...any) {
	if !tee.enabled(LevelAll) {
		return
	}
	fmt.Fprintln(tee.GetWriter(LevelAll), a...)
}

func (tee *TeeLogger) Printlnf(format string, a ...any) {
	if !tee.enabled(LevelAll) {
		return
	}
	w := tee.GetWriter(LevelAll)
	fmt.Fprintf(w, format, a...)
	io.WriteString(w, "\n")
}

func (tee *TeeLogger) output(calldepth int, level Level, msg string) error {
	var err error
	for _, l := range tee.loggers {
		if !l.enabled(level) {
			continue
		}
		err = errors.Join(err, l.output(calldepth+1, level, msg))
	}

	return err
}

func (tee *TeeLogger) Debug(a ...any) {
	if !tee.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprint(a...)
	tee.output(3, LevelDebug, msg)
}

func (tee *TeeLogger) Debugf(format string, a ...any) {
	if !tee.enabled(LevelDebug) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	tee.output(3, LevelDebug, msg)
}

func (tee *TeeLogger) Info(a ...any) {
	if !tee.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprint(a...)
	tee.output(3, LevelInfo, msg)
}

func (tee *TeeLogger) Infof(format string, a ...any) {
	if !tee.enabled(LevelInfo) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	tee.output(3, LevelInfo, msg)
}

func (tee *TeeLogger) Warn(a ...any) {
	if !tee.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprint(a...)
	tee.output(3, LevelWarn, msg)
}

func (tee *TeeLogger) Warnf(format string, a ...any) {
	if !tee.enabled(LevelWarn) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	tee.output(3, LevelWarn, msg)
}

func (tee *TeeLogger) Error(a ...any) {
	if !tee.enabled(LevelError) {
		return
	}
	msg := fmt.Sprint(a...)
	tee.output(3, LevelError, msg)
}

func (tee *TeeLogger) Errorf(format string, a ...any) {
	if !tee.enabled(LevelError) {
		return
	}
	msg := fmt.Sprintf(format, a...)
	tee.output(3, LevelError, msg)
}

func (tee *TeeLogger) ErrorAt(err error, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprint(a...))
	if !tee.enabled(LevelError) {
		return err
	}
	tee.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func (tee *TeeLogger) ErrorAtf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	err = errors.WithStack(err, 4, fmt.Sprintf(format, a...))
	if !tee.enabled(LevelError) {
		return err
	}
	tee.output(3, LevelError, errors.GetErrorDetails(err))
	return err
}

func (tee *TeeLogger) Fatal(a ...any) {
	if tee.enabled(LevelFatal) {
		msg := fmt.Sprint(a...)
		tee.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func (tee *TeeLogger) Fatalf(format string, a ...any) {
	if tee.enabled(LevelFatal) {
		msg := fmt.Sprintf(format, a...)
		tee.output(3, LevelFatal, msg)
	}
	os.Exit(1)
}

func (tee *TeeLogger) Panic(a ...any) {
	msg := fmt.Sprint(a...)
	if tee.enabled(LevelPanic) {
		tee.output(3, LevelPanic, msg)
	}
	panic(msg)
}

func (tee *TeeLogger) Panicf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if tee.enabled(LevelPanic) {
		tee.output(3, LevelPanic, msg)
	}
	panic(msg)
}

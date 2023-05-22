package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
)

const (
	defaultOmmittedErrMessage = "something is wrong"
)

type stackErr struct {
	msg   string
	cause error
	stack stack
}

func Error(a ...any) error {
	return WithStack(nil, 4, fmt.Sprint(a...))
}

func Errorf(format string, a ...any) error {
	return WithStack(nil, 4, fmt.Sprintf(format, a...))
}

func ErrorAt(err error, a ...any) error {
	return WithStack(err, 4, fmt.Sprint(a...))
}

func ErrorAtf(err error, format string, a ...any) error {
	return WithStack(err, 4, fmt.Sprintf(format, a...))
}

func WithStack(err error, callDepth int, msg string) error {
	werr := &stackErr{}
	if len(msg) == 0 {
		msg = defaultOmmittedErrMessage
	}
	werr.msg = msg
	werr.stack = getStack(callDepth)
	werr.cause = err
	return werr
}

func RootCause(err error) error {
	if err == nil {
		return nil
	}

	for {
		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return err
		}
		next := u.Unwrap()
		if next == nil {
			return err
		}
		err = next
	}
}

func GetErrorDetails(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("\n%+v", err)
}

type stack []runtime.Frame

func getStack(skip int) stack {
	callers := make([]uintptr, 64)
	n := runtime.Callers(skip, callers)

	s := make([]runtime.Frame, 0, 5)
	frames := runtime.CallersFrames(callers[:n])
	for {
		frame, more := frames.Next()
		s = append(s, frame)
		if !more {
			break
		}
	}

	return s
}

func (e *stackErr) Error() string {
	return e.msg
}

func (e *stackErr) Unwrap() error {
	return e.cause
}

func (e *stackErr) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if e.cause != nil {
			if cause, ok := e.cause.(fmt.Formatter); ok {
				cause.Format(state, verb)
			} else {
				fmt.Fprintf(state, "* %v\n", e.cause)
			}
		}

		fmt.Fprintf(state, "* %s\n", e.msg)
		if state.Flag('+') {
			frameCount, showFrame := state.Width()
			if !showFrame {
				frameCount = 1
			}
			frames := e.stack
			if len(frames) > frameCount {
				frames = e.stack[:frameCount]
			}
			for _, frame := range frames {
				fmt.Fprintf(state, "\t* %s:%d %s\n", frame.File, frame.Line, path.Base(frame.Function))
			}
		}
	case 's':
		io.WriteString(state, e.msg)
	case 'q':
		fmt.Fprintf(state, "%q", e.msg)
	}
}

package errors

import (
	"fmt"
	"io"
	"strings"
)

type multiErr struct {
	errs []error
}

func Join(errs ...error) error {
	n := 0
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}
	if len(nonNilErrs) == 0 {
		return nil
	}
	e := &multiErr{
		errs: make([]error, 0, n),
	}

	if e0, ok := nonNilErrs[0].(*multiErr); ok {
		e = e0
	} else {
		e.errs = append(e.errs, nonNilErrs[0])
	}

	e.errs = append(e.errs, nonNilErrs[1:]...)

	return e
}

func (e *multiErr) Error() string {
	sb := &strings.Builder{}
	for i, err := range e.errs {
		sb.WriteString(err.Error())
		if i != len(e.errs)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (e *multiErr) Unwrap() []error {
	return e.errs
}

func (e *multiErr) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		for i, err := range e.errs {
			fmt.Fprintf(state, "%d. %s\n", i+1, RootCause(err).Error())
			if _, ok := err.(fmt.Formatter); ok {
				sb := &strings.Builder{}
				format := revertFormatState(state, verb)
				fmt.Fprintf(sb, format, err)
				lines := strings.Split(sb.String(), "\n")
				for _, line := range lines {
					if line != "" {
						fmt.Fprintf(state, "\t%s\n", line)
					}
				}
			}
		}
	case 's':
		for _, err := range e.errs {
			fmt.Fprintf(state, "%s\n", err.Error())
		}
	case 'q':
		for _, err := range e.errs {
			fmt.Fprintf(state, "%q\n", err.Error())
		}
	}
}

func revertFormatState(state fmt.State, verb rune) string {
	sb := &strings.Builder{}
	io.WriteString(sb, "%")
	if state.Flag('+') {
		io.WriteString(sb, "+")
	}
	if state.Flag('#') {
		io.WriteString(sb, "#")
	}

	if wid, ok := state.Width(); ok {
		fmt.Fprintf(sb, "%d", wid)
	}

	if prec, ok := state.Precision(); ok {
		fmt.Fprintf(sb, ".%d", prec)
	}

	fmt.Fprintf(sb, "%c", verb)

	return sb.String()
}

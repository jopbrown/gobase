//go:build go1.13
// +build go1.13

package errors

import (
	stderrors "errors"
)

func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target any) bool {
	return stderrors.As(err, target)
}

func AsIs[E error](err error) (E, bool) {
	var target E
	ok := stderrors.As(err, &target)
	return target, ok
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

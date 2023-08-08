//go:build go1.18
// +build go1.18

package errors_test

import (
	"io"
	"os"
	"testing"

	"github.com/jopbrown/gobase/errors"
	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	err := errors.ErrorAt(io.EOF)
	err = errors.ErrorAt(err)
	err = errors.Join(err, os.ErrClosed)
	err = errors.ErrorAt(err)

	assert.True(t, errors.Is(err, io.EOF))
	assert.True(t, errors.Is(err, os.ErrClosed))
}

type CustomError1 string

func (err CustomError1) Error() string {
	return string(err)
}

type CustomError2 string

func (err CustomError2) Error() string {
	return string(err)
}

func TestAs(t *testing.T) {
	err := errors.ErrorAt(CustomError1("err1"))
	err = errors.ErrorAt(err)
	err = errors.Join(err, CustomError2("err2"))
	err = errors.ErrorAt(err)

	var err1 CustomError1
	assert.True(t, errors.As(err, &err1))
	assert.Equal(t, "err1", err1.Error())

	var err2 CustomError2
	assert.True(t, errors.As(err, &err2))
	assert.Equal(t, "err2", err2.Error())

	var ok bool
	err1, ok = errors.AsIs[CustomError1](err)
	assert.True(t, ok)
	assert.Equal(t, "err1", err1.Error())

	err2, ok = errors.AsIs[CustomError2](err)
	assert.True(t, ok)
	assert.Equal(t, "err1", err1.Error())
}

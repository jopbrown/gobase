package errors_test

import (
	"testing"

	"github.com/jopbrown/gobase/errors"
	"github.com/stretchr/testify/assert"
)

var baseErr = errors.Error("an error")

func Err() error {
	return baseErr
}

func Err1() (int, error) {
	return 1, baseErr
}

func Err2() (int, int, error) {
	return 1, 2, baseErr
}

func Err3() (int, int, int, error) {
	return 1, 2, 3, baseErr
}

func NoErr() error {
	return nil
}

func NoErr1() (int, error) {
	return 1, nil
}

func NoErr2() (int, int, error) {
	return 1, 2, nil
}

func NoErr3() (int, int, int, error) {
	return 1, 2, 3, nil
}

func tuple2Slice2[T1, T2 any](v1 T1, v2 T2) []any {
	return []any{v1, v2}
}

func tuple2Slice3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3) []any {
	return []any{v1, v2, v3}
}

func tuple2Slice4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4) []any {
	return []any{v1, v2, v3, v4}
}

func TestMust(t *testing.T) {
	assert.Panics(t, func() { errors.Must(Err()) })
	assert.Panics(t, func() { errors.Must1(Err1()) })
	assert.Panics(t, func() { errors.Must2(Err2()) })
	assert.Panics(t, func() { errors.Must3(Err3()) })

	assert.NotPanics(t, func() { errors.Must(NoErr()) })
	assert.NotPanics(t, func() { errors.Must1(NoErr1()) })
	assert.NotPanics(t, func() { errors.Must2(NoErr2()) })
	assert.NotPanics(t, func() { errors.Must3(NoErr3()) })
}

func TestShould(t *testing.T) {
	assert.Equal(t, []any{1, true}, tuple2Slice2(errors.Should1(Err1())))
	assert.Equal(t, []any{1, 2, true}, tuple2Slice3(errors.Should2(Err2())))
	assert.Equal(t, []any{1, 2, 3, true}, tuple2Slice4(errors.Should3(Err3())))

	assert.Equal(t, []any{1, false}, tuple2Slice2(errors.Should1(NoErr1())))
	assert.Equal(t, []any{1, 2, false}, tuple2Slice3(errors.Should2(NoErr2())))
	assert.Equal(t, []any{1, 2, 3, false}, tuple2Slice4(errors.Should3(NoErr3())))
}

func TestHas(t *testing.T) {
	assert.True(t, errors.Has(Err()))
	assert.True(t, errors.Has1(Err1()))
	assert.True(t, errors.Has2(Err2()))
	assert.True(t, errors.Has3(Err3()))

	assert.False(t, errors.Has(NoErr()))
	assert.False(t, errors.Has1(NoErr1()))
	assert.False(t, errors.Has2(NoErr2()))
	assert.False(t, errors.Has3(NoErr3()))
}

func TestIgnore(t *testing.T) {
	assert.Equal(t, 1, errors.Ignore1(Err1()))
	assert.Equal(t, []any{1, 2}, tuple2Slice2(errors.Ignore2(Err2())))
	assert.Equal(t, []any{1, 2, 3}, tuple2Slice3(errors.Ignore3(Err3())))

	assert.Equal(t, 1, errors.Ignore1(NoErr1()))
	assert.Equal(t, []any{1, 2}, tuple2Slice2(errors.Ignore2(NoErr2())))
	assert.Equal(t, []any{1, 2, 3}, tuple2Slice3(errors.Ignore3(NoErr3())))
}

func TestGet(t *testing.T) {
	assert.Error(t, errors.Get1(Err1()))
	assert.Error(t, errors.Get2(Err2()))
	assert.Error(t, errors.Get3(Err3()))

	assert.NoError(t, errors.Get1(NoErr1()))
	assert.NoError(t, errors.Get2(NoErr2()))
	assert.NoError(t, errors.Get3(NoErr3()))
}

func TestCatch(t *testing.T) {
	assert.Equal(t, baseErr, errors.Catch(func() {
		panic(baseErr)
	}))

	assert.Equal(t, nil, errors.Catch(func() {
		// no panic
	}))

	assert.Equal(t, "anything", errors.Catch(func() {
		panic("anything")
	}).Error())
}

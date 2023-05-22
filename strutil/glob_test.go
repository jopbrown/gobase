package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustComplieGlob_Asterisk(t *testing.T) {
	matcher := MustComplieGlob("a*b.txt")
	assert.True(t, matcher.MatchString(`ab.txt`))
	assert.True(t, matcher.MatchString(`aabb.txt`))
	assert.True(t, matcher.MatchString(`a**b.txt`))
	assert.True(t, matcher.MatchString(`axxxb.txt`))
	assert.False(t, matcher.MatchString(`abc.txt`))
	assert.False(t, matcher.MatchString(`axb.dat`))
}

func TestMustComplieGlob_Question(t *testing.T) {
	matcher := MustComplieGlob("a?b.txt")
	assert.False(t, matcher.MatchString(`ab.txt`))
	assert.False(t, matcher.MatchString(`aabb.txt`))
	assert.False(t, matcher.MatchString(`a**b.txt`))
	assert.True(t, matcher.MatchString(`axb.txt`))
	assert.False(t, matcher.MatchString(`abc.txt`))
	assert.False(t, matcher.MatchString(`axb.dat`))
}

func TestMustComplieGlob_Mix(t *testing.T) {
	matcher := MustComplieGlob("a*b?c.txt")
	assert.False(t, matcher.MatchString(`abc.txt`))
	assert.True(t, matcher.MatchString(`abxc.txt`))
	assert.True(t, matcher.MatchString(`axbxc.txt`))
}

func TestMustComplieGlob_None(t *testing.T) {
	matcher := MustComplieGlob("abc.txt")
	assert.True(t, matcher.MatchString(`abc.txt`))
	assert.False(t, matcher.MatchString(`abxc.txt`))
	assert.False(t, matcher.MatchString(`axbxc.txt`))
}

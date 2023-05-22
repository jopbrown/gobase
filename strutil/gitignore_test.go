package strutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule1_EmptyLine(t *testing.T) {
	matcher := MustCompileGitIgnore("")
	assert.False(t, matcher.MatchString("hello"))
	assert.False(t, matcher.MatchString("/world"))
	assert.False(t, matcher.MatchString("any/path"))
	assert.False(t, matcher.MatchString("/any/path"))
	assert.False(t, matcher.MatchString("/"))
	assert.False(t, matcher.MatchString(""))
}

func TestRule2_Comment(t *testing.T) {
	matcher := MustCompileGitIgnore("#ping")
	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("#ping"))
	assert.False(t, matcher.MatchString("any/path/#ping"))
	assert.False(t, matcher.MatchString("/any/path/#ping"))
	assert.False(t, matcher.MatchString("#ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/#ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/#ping/z/y/x"))
}

func TestRule2_EscComment(t *testing.T) {
	matcher := MustCompileGitIgnore(`\#ping`)
	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("#ping"))
	assert.True(t, matcher.MatchString("any/path/#ping"))
	assert.True(t, matcher.MatchString("/any/path/#ping"))
	assert.True(t, matcher.MatchString("#ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/#ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/#ping/z/y/x"))
}

func TestRule4_Negate(t *testing.T) {
	matcher := MustCompileGitIgnore("!ping")
	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("!ping"))
	assert.False(t, matcher.MatchString("any/path/!ping"))
	assert.False(t, matcher.MatchString("/any/path/!ping"))
	assert.False(t, matcher.MatchString("!ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/!ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/!ping/z/y/x"))

	assert.False(t, matcher.MatchString("pong"))
	assert.False(t, matcher.MatchString("any/path/pong"))
	assert.False(t, matcher.MatchString("/any/path/pong"))
	assert.False(t, matcher.MatchString("!pong"))
	assert.False(t, matcher.MatchString("any/path/!pong"))
	assert.False(t, matcher.MatchString("/any/path/!pong"))
	assert.False(t, matcher.MatchString("!pong/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/!pong/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/!pong/z/y/x"))
}

func TestRule4_EscNegate(t *testing.T) {
	matcher := MustCompileGitIgnore(`\!ping`)
	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("!ping"))
	assert.True(t, matcher.MatchString("any/path/!ping"))
	assert.True(t, matcher.MatchString("/any/path/!ping"))
	assert.True(t, matcher.MatchString("!ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/!ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/!ping/z/y/x"))
}

func TestRule6_NoSlash(t *testing.T) {
	matcher := MustCompileGitIgnore(`ping`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping"))
	assert.False(t, matcher.MatchString("/head_ping"))
	assert.False(t, matcher.MatchString("head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping"))
	assert.False(t, matcher.MatchString("/any/path/head_ping"))
	assert.False(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.False(t, matcher.MatchString("ping_foot"))
	assert.False(t, matcher.MatchString("/ping_foot"))
	assert.False(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot"))
	assert.False(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule6_WithSlash_Leading(t *testing.T) {
	matcher := MustCompileGitIgnore(`/ping`)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))

	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))

	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping"))
	assert.False(t, matcher.MatchString("/head_ping"))
	assert.False(t, matcher.MatchString("head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping"))
	assert.False(t, matcher.MatchString("/any/path/head_ping"))
	assert.False(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.False(t, matcher.MatchString("ping_foot"))
	assert.False(t, matcher.MatchString("/ping_foot"))
	assert.False(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot"))
	assert.False(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule6_WithSlash_Middle(t *testing.T) {
	matcher := MustCompileGitIgnore(`any/path/ping`)

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))

	assert.False(t, matcher.MatchString("ping/z/y/x"))
	assert.False(t, matcher.MatchString("/ping/z/y/x"))

	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))

	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))
}

func TestRule7_Simple(t *testing.T) {
	require.NoError(t, os.RemoveAll("tmp"))
	require.NoError(t, os.MkdirAll("tmp", 0755))
	require.NoError(t, os.Chdir("tmp"))
	defer os.Chdir("..")

	matcher := MustCompileGitIgnore(`ping/`)

	os.RemoveAll("ping")

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	os.MkdirAll("ping", 0755)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	os.RemoveAll("any/path/ping")

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	os.MkdirAll("any/path/ping", 0755)

	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule7_MixRule6(t *testing.T) {
	require.NoError(t, os.RemoveAll("tmp"))
	require.NoError(t, os.MkdirAll("tmp", 0755))
	require.NoError(t, os.Chdir("tmp"))
	defer os.Chdir("..")

	matcher := MustCompileGitIgnore(`/ping/`)

	os.RemoveAll("ping")

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	os.MkdirAll("ping", 0755)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	os.RemoveAll("any/path/ping")

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))

	os.MkdirAll("any/path/ping", 0755)

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule9_2_Trail(t *testing.T) {
	matcher := MustCompileGitIgnore(`ping*`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping"))
	assert.False(t, matcher.MatchString("/head_ping"))
	assert.False(t, matcher.MatchString("head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping"))
	assert.False(t, matcher.MatchString("/any/path/head_ping"))
	assert.False(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.True(t, matcher.MatchString("ping_foot"))
	assert.True(t, matcher.MatchString("/ping_foot"))
	assert.True(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping_foot"))
	assert.True(t, matcher.MatchString("/any/path/ping_foot"))
	assert.True(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule9_2_Lead(t *testing.T) {
	matcher := MustCompileGitIgnore(`*ping`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.True(t, matcher.MatchString("head_ping"))
	assert.True(t, matcher.MatchString("/head_ping"))
	assert.True(t, matcher.MatchString("head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/head_ping"))
	assert.True(t, matcher.MatchString("/any/path/head_ping"))
	assert.True(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.False(t, matcher.MatchString("ping_foot"))
	assert.False(t, matcher.MatchString("/ping_foot"))
	assert.False(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot"))
	assert.False(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule9_2_LeadAndTrail(t *testing.T) {
	matcher := MustCompileGitIgnore(`*ping*`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.True(t, matcher.MatchString("head_ping"))
	assert.True(t, matcher.MatchString("/head_ping"))
	assert.True(t, matcher.MatchString("head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/head_ping"))
	assert.True(t, matcher.MatchString("/any/path/head_ping"))
	assert.True(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.True(t, matcher.MatchString("ping_foot"))
	assert.True(t, matcher.MatchString("/ping_foot"))
	assert.True(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping_foot"))
	assert.True(t, matcher.MatchString("/any/path/ping_foot"))
	assert.True(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.True(t, matcher.MatchString("head_ping_foot"))
	assert.True(t, matcher.MatchString("/head_ping_foot"))
	assert.True(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.True(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.True(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule9_2_Middle(t *testing.T) {
	matcher := MustCompileGitIgnore(`p*g`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping"))
	assert.False(t, matcher.MatchString("/head_ping"))
	assert.False(t, matcher.MatchString("head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping"))
	assert.False(t, matcher.MatchString("/any/path/head_ping"))
	assert.False(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.False(t, matcher.MatchString("ping_foot"))
	assert.False(t, matcher.MatchString("/ping_foot"))
	assert.False(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot"))
	assert.False(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule9_2_LeadAndMiddle(t *testing.T) {
	matcher := MustCompileGitIgnore(`*p*g`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.True(t, matcher.MatchString("head_ping"))
	assert.True(t, matcher.MatchString("/head_ping"))
	assert.True(t, matcher.MatchString("head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/head_ping"))
	assert.True(t, matcher.MatchString("/any/path/head_ping"))
	assert.True(t, matcher.MatchString("any/path/head_ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/head_ping/z/y/x"))

	assert.False(t, matcher.MatchString("ping_foot"))
	assert.False(t, matcher.MatchString("/ping_foot"))
	assert.False(t, matcher.MatchString("ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot"))
	assert.False(t, matcher.MatchString("any/path/ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping_foot/z/y/x"))

	assert.False(t, matcher.MatchString("head_ping_foot"))
	assert.False(t, matcher.MatchString("/head_ping_foot"))
	assert.False(t, matcher.MatchString("head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot"))
	assert.False(t, matcher.MatchString("any/path/head_ping_foot/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/head_ping_foot/z/y/x"))
}

func TestRule9_1_Middle(t *testing.T) {
	matcher := MustCompileGitIgnore(`pi?g`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("pig"))
	assert.False(t, matcher.MatchString("/pig"))
	assert.False(t, matcher.MatchString("pig/z/y/x"))
	assert.False(t, matcher.MatchString("/pig/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/pig"))
	assert.False(t, matcher.MatchString("/any/path/pig"))
	assert.False(t, matcher.MatchString("any/path/pig/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/pig/z/y/x"))
}

func TestRule9_1_Middle_Omit(t *testing.T) {
	matcher := MustCompileGitIgnore(`pin?g`)
	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))
	assert.False(t, matcher.MatchString("ping/z/y/x"))
	assert.False(t, matcher.MatchString("/ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule9_2_Middle_Omit(t *testing.T) {
	matcher := MustCompileGitIgnore(`pin*g`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule9_2_Middle_MixRule6(t *testing.T) {
	matcher := MustCompileGitIgnore(`/p*g`)
	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule9_2_Middle_MixRule7(t *testing.T) {
	require.NoError(t, os.RemoveAll("tmp"))
	require.NoError(t, os.MkdirAll("tmp", 0755))
	require.NoError(t, os.Chdir("tmp"))
	defer os.Chdir("..")

	matcher := MustCompileGitIgnore(`p*g/`)

	os.RemoveAll("ping")

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	os.MkdirAll("ping", 0755)
	os.MkdirAll("any/path/ping", 0755)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))
	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))
}

func TestRule10_1_Leading(t *testing.T) {
	matcher := MustCompileGitIgnore(`**/ping`)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))

	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))

	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.True(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.True(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.True(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))
}

func TestRule10_2_Trailing(t *testing.T) {
	require.NoError(t, os.RemoveAll("tmp"))
	require.NoError(t, os.MkdirAll("tmp", 0755))
	require.NoError(t, os.Chdir("tmp"))
	defer os.Chdir("..")

	matcher := MustCompileGitIgnore(`ping/**`)

	os.RemoveAll("ping")

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))

	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))

	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))

	os.MkdirAll("ping", 0755)
	os.MkdirAll("any/path/ping", 0755)
	os.MkdirAll("/a/b/c/any/path/ping", 0755)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))

	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))

	assert.False(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))
}

func TestRule10_3_LeadingAndTrailing(t *testing.T) {
	matcher := MustCompileGitIgnore(`any/**/ping`)

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))

	assert.False(t, matcher.MatchString("ping/z/y/x"))
	assert.False(t, matcher.MatchString("/ping/z/y/x"))

	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))

	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))
}

func TestRule10_4_Other(t *testing.T) {
	matcher := MustCompileGitIgnore(`p**ng`)

	assert.True(t, matcher.MatchString("ping"))
	assert.True(t, matcher.MatchString("/ping"))

	assert.True(t, matcher.MatchString("ping/z/y/x"))
	assert.True(t, matcher.MatchString("/ping/z/y/x"))

	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))

	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.True(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.True(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.True(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))
}

func TestRule10_MixOtherRules(t *testing.T) {
	require.NoError(t, os.RemoveAll("tmp"))
	require.NoError(t, os.MkdirAll("tmp", 0755))
	require.NoError(t, os.Chdir("tmp"))
	defer os.Chdir("..")

	matcher := MustCompileGitIgnore(`a?y/**/p*g/`)

	assert.False(t, matcher.MatchString("ping"))
	assert.False(t, matcher.MatchString("/ping"))

	assert.False(t, matcher.MatchString("ping/z/y/x"))
	assert.False(t, matcher.MatchString("/ping/z/y/x"))

	assert.False(t, matcher.MatchString("any/path/ping"))
	assert.False(t, matcher.MatchString("/any/path/ping"))

	assert.True(t, matcher.MatchString("any/path/ping/z/y/x"))
	assert.True(t, matcher.MatchString("/any/path/ping/z/y/x"))

	assert.False(t, matcher.MatchString("a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping"))
	assert.False(t, matcher.MatchString("a/b/c/any/path/ping/z/y/x"))
	assert.False(t, matcher.MatchString("/a/b/c/any/path/ping/z/y/x"))

	os.MkdirAll("any/path/ping", 0755)
	assert.True(t, matcher.MatchString("any/path/ping"))
	assert.True(t, matcher.MatchString("/any/path/ping"))
}

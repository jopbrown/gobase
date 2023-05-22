package fsutil

import (
	"path/filepath"
	"testing"

	"github.com/jopbrown/gobase/strutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListWithMatcher(t *testing.T) {
	actual, err := ListWithMatcher(".", strutil.MustComplieGlob("*.go"))
	require.NoError(t, err)

	golden, err := filepath.Glob("*.go")
	require.NoError(t, err)

	assert.Equal(t, golden, actual)
}

package rotate_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jopbrown/gobase/log/rotate"
	"github.com/jopbrown/gobase/must"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenFile(t *testing.T) {
	os.RemoveAll("tmp")

	w, err := rotate.OpenFile("tmp/test.log", 300*time.Millisecond, 0)
	require.NoError(t, err)

	fmt.Fprintf(w, "aaa")
	fmt.Fprintf(w, "bbb")

	time.Sleep(300 * time.Millisecond)

	fmt.Fprintf(w, "ccc")
	fmt.Fprintf(w, "ddd")

	time.Sleep(300 * time.Millisecond)

	fmt.Fprintf(w, "eee")
	fmt.Fprintf(w, "fff")

	assert.Equal(t, "aaabbbccc", string(must.Value(os.ReadFile(must.Value(filepath.Glob("tmp/*_01.log"))[0]))))
	assert.Equal(t, "dddeee", string(must.Value(os.ReadFile(must.Value(filepath.Glob("tmp/*_02.log"))[0]))))
	assert.Equal(t, "fff", string(must.Value(os.ReadFile("tmp/test.log"))))

	w.Close()
}

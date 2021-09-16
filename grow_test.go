package mmap

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ImSingee/tt"
)

func TestDeleteFileBeforeGrow(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadWrite(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	err = os.RemoveAll(filepath.Dir(f))
	tt.AssertIsNotError(t, err)

	n, err := mmap.WriteAt([]byte{6}, 128) // will grow
	tt.AssertIsError(t, err)
	tt.AssertTrue(t, errors.Is(err, os.ErrNotExist))
	tt.AssertEqual(t, 0, n)
	tt.AssertTrue(t, mmap.IsClosed())
}

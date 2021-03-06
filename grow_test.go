package mmap

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ImSingee/tt"
)

func TestGrow(t *testing.T) {
	t.Run("grow-to-twice", func(t *testing.T) {
		mmap, err := New(NewReadWrite(""))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		tt.AssertEqual(t, oneMB, mmap.Cap())

		err = mmap.EnsureCapacity(4 * oneMB) // grow to 4 MB for test
		tt.AssertEqual(t, 4*oneMB, mmap.Cap())

		// won't grow if write to right bound
		_, err = mmap.WriteAt([]byte{1}, 4*oneMB-1)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, 4*oneMB, mmap.Cap())

		// will grow to twice (8MB) if write something to [right bound + 1]
		_, err = mmap.WriteAt([]byte{2}, 4*oneMB)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, 8*oneMB, mmap.Cap())
	})

	t.Run("grow-to-large", func(t *testing.T) {
		mmap, err := New(NewReadWrite(""))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		tt.AssertEqual(t, oneMB, mmap.Cap())

		// if current is not less than 2GB, simple add one GB (not double)
		err = mmap.EnsureCapacity(2 * oneGB) // grow to 2 GB for test
		tt.AssertEqual(t, 2*oneGB, mmap.Cap())

		_, err = mmap.WriteAt([]byte{1}, 2*oneGB)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, 3*oneGB, mmap.Cap())
	})

	t.Run("grow-align-to-1MB", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		tt.AssertEqual(t, LenOfHelloWorld, mmap.Cap()) // < 512KB

		// write anything, it will grow to 1MB (because of align)
		_, err = mmap.WriteAt([]byte{1}, 16)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, oneMB, mmap.Cap())
	})

	t.Run("grow-to-need", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		tt.AssertEqual(t, LenOfHelloWorld, mmap.Cap()) // < 1MB

		// write anything at 2MB, it will grow to 3MB (because of need and align)
		_, err = mmap.WriteAt([]byte{1}, 2*oneMB)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, 3*oneMB, mmap.Cap())
	})
}

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

func TestGrowAfterClose(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadWrite(f))
	tt.AssertIsNotError(t, err)
	closeMmap(t, mmap)

	n, err := mmap.WriteAt([]byte{6}, 1024)
	tt.AssertEqual(t, ErrIsClosed, err)
	tt.AssertEqual(t, 0, n)
}

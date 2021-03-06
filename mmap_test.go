package mmap

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ImSingee/tt"
)

const HelloWorld = "Hello world!"
const LenOfHelloWorld = len(HelloWorld)

func TestMmapNewReadOnly(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	tt.AssertEqual(t, LenOfHelloWorld, mmap.Cap())
	tt.AssertEqual(t, LenOfHelloWorld, len(mmap.data))
	tt.AssertEqual(t, LenOfHelloWorld, cap(mmap.data))
	tt.AssertEqual(t, []byte(HelloWorld), mmap.data)
}

func TestMmapNewEmptyReadWrite(t *testing.T) {
	mmap, err := New(NewReadWrite(""))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	f := mmap.args.(*Args).File

	tt.AssertEqual(t, DefaultInitLength, mmap.Cap())
	tt.AssertEqual(t, DefaultInitLength, len(mmap.data))
	tt.AssertEqual(t, DefaultInitLength, cap(mmap.data))
	tt.AssertEqual(t, make([]byte, DefaultInitLength), mmap.data)
	tt.AssertEqual(t, int64(DefaultInitLength), fileSize(f))

	n, err := mmap.WriteAt([]byte(HelloWorld), 8)
	tt.AssertEqual(t, LenOfHelloWorld, n)
	tt.AssertIsNotError(t, err)

	tt.AssertEqual(t, DefaultInitLength, mmap.Cap())
	tt.AssertEqual(t, DefaultInitLength, len(mmap.data))
	tt.AssertEqual(t, DefaultInitLength, cap(mmap.data))

	expect := make([]byte, DefaultInitLength)
	copy(expect[8:], []byte(HelloWorld))
	tt.AssertEqual(t, expect, mmap.data)

	p, err := ioutil.ReadFile(f)
	tt.AssertIsNotError(t, err)
	tt.AssertEqual(t, expect, p)
}

func TestMmapNewReadWrite(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadWrite(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	tt.AssertEqual(t, LenOfHelloWorld, mmap.Cap())
	tt.AssertEqual(t, LenOfHelloWorld, len(mmap.data))
	tt.AssertEqual(t, LenOfHelloWorld, cap(mmap.data))
	tt.AssertEqual(t, []byte(HelloWorld), mmap.data)

	tt.AssertEqual(t, int64(LenOfHelloWorld), fileSize(f))

	n, err := mmap.WriteAt([]byte(HelloWorld), int64(LenOfHelloWorld))
	tt.AssertEqual(t, LenOfHelloWorld, n)
	tt.AssertIsNotError(t, err)

	tt.AssertEqual(t, oneMB, mmap.Cap())
	tt.AssertEqual(t, oneMB, len(mmap.data))
	tt.AssertEqual(t, oneMB, cap(mmap.data))

	expect := make([]byte, oneMB)
	copy(expect, []byte(HelloWorld))
	copy(expect[LenOfHelloWorld:], []byte(HelloWorld))
	tt.AssertEqual(t, expect, mmap.data)

	p, err := ioutil.ReadFile(f)
	tt.AssertIsNotError(t, err)
	tt.AssertEqual(t, expect, p)
}

func TestNotExistFile(t *testing.T) {
	t.Run("not-exist-readonly", func(t *testing.T) {
		mmap, err := New(&Args{
			File:       "/path/to/not-exist",
			InitLength: 1,
			Readonly:   true,
			Private:    false,
		})
		tt.AssertIsError(t, err)
		tt.AssertTrue(t, os.IsNotExist(err))
		tt.AssertTrue(t, errors.Is(err, os.ErrNotExist))

		tt.AssertIsNotNil(t, mmap)
		tt.AssertTrue(t, mmap.IsClosed())
	})

	t.Run("not-exist-readwrite", func(t *testing.T) {
		mmap, err := New(&Args{
			File:       "/path/to/not-exist",
			InitLength: 1,
			Readonly:   false,
			Private:    false,
		})
		tt.AssertIsError(t, err)
		tt.AssertTrue(t, os.IsNotExist(err))
		tt.AssertTrue(t, errors.Is(err, os.ErrNotExist))

		tt.AssertIsNotNil(t, mmap)
		tt.AssertTrue(t, mmap.IsClosed())
	})
}

func closeMmap(t *testing.T, m *Mmap) {
	t.Helper()

	tt.AssertIsNotError(t, m.Close())
}

func newHelloWorldFile() (string, error) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(d, "")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(HelloWorld)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func fileSize(f string) int64 {
	stat, err := os.Stat(f)
	if err != nil {
		return -1
	}
	return stat.Size()
}

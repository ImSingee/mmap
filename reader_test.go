package mmap

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ImSingee/tt"
)

func TestRead(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadOnly(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		p := make([]byte, 4)
		n, err := mmap.ReadAt(p, 1)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, 4, n)
		tt.AssertEqual(t, "ello", string(p))
	})

	t.Run("read more", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadOnly(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		p := make([]byte, 128)
		n, err := mmap.ReadAt(p, 1)
		tt.AssertEqual(t, io.EOF, err)
		tt.AssertEqual(t, LenOfHelloWorld-1, n)

		expect := make([]byte, 128)
		copy(expect, HelloWorld[1:])

		tt.AssertEqual(t, expect, p)
	})

	t.Run("read not exist", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadOnly(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		p := make([]byte, 128)
		n, err := mmap.ReadAt(p, 666)
		tt.AssertEqual(t, io.EOF, err)
		tt.AssertEqual(t, 0, n)
		tt.AssertEqual(t, make([]byte, 128), p)
	})
}

func TestBytes(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	p, err := mmap.Bytes(1, 4)
	tt.AssertIsNotError(t, err)
	tt.AssertEqual(t, "ello", string(p))
}

func TestBytesMore(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	p, err := mmap.Bytes(1, 128)
	tt.AssertEqual(t, io.EOF, err)
	tt.AssertEqual(t, HelloWorld[1:], string(p))
}

func TestBytesNotExist(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	p, err := mmap.Bytes(666, 128)
	tt.AssertEqual(t, io.EOF, err)
	tt.AssertEqual(t, 0, len(p))
}

func TestWriteTo(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	buf := &bytes.Buffer{}

	n, err := mmap.WriteTo(buf)
	tt.AssertIsNotError(t, err)
	tt.AssertEqual(t, int64(LenOfHelloWorld), n)
	tt.AssertEqual(t, HelloWorld, buf.String())
}

func TestWriteToAt(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	buf := &bytes.Buffer{}

	n, err := mmap.WriteToAt(6, buf)
	tt.AssertIsNotError(t, err)
	tt.AssertEqual(t, int64(LenOfHelloWorld)-6, n)
	tt.AssertEqual(t, HelloWorld[6:], buf.String())
}

func TestWriteToAtNotExist(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	defer closeMmap(t, mmap)

	buf := &bytes.Buffer{}
	n, err := mmap.WriteToAt(666, buf)
	tt.AssertEqual(t, io.EOF, err)
	tt.AssertEqual(t, int64(0), n)
	tt.AssertEqual(t, 0, buf.Len())
}

func TestReadAfterClose(t *testing.T) {
	f, err := newHelloWorldFile()
	tt.AssertIsNotError(t, err)

	mmap, err := New(NewReadOnly(f))
	tt.AssertIsNotError(t, err)
	closeMmap(t, mmap)

	n, err := mmap.ReadAt(nil, 1)
	tt.AssertEqual(t, ErrIsClosed, err)
	tt.AssertEqual(t, 0, n)

	p, err := mmap.Bytes(0, 3)
	tt.AssertEqual(t, ErrIsClosed, err)
	tt.AssertEqual(t, 0, len(p))

	buf := new(bytes.Buffer)
	n_, err := mmap.WriteTo(buf)
	tt.AssertEqual(t, ErrIsClosed, err)
	tt.AssertEqual(t, int64(0), n_)
	tt.AssertEqual(t, 0, buf.Len())

	buf = new(bytes.Buffer)
	n_, err = mmap.WriteToAt(6, buf)
	tt.AssertEqual(t, ErrIsClosed, err)
	tt.AssertEqual(t, int64(0), n_)
	tt.AssertEqual(t, 0, buf.Len())
}

func TestReader(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadOnly(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		r := mmap.ReaderAt(1)
		all, err := io.ReadAll(r)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, HelloWorld[1:], string(all))
	})

	t.Run("read-multiple-times", func(t *testing.T) {
		mmap, err := New(NewReadWrite(""))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		data := []byte(strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1024))

		_, err = mmap.WriteAt(data, 0)
		tt.AssertIsNotError(t, err)

		buf := make([]byte, 13)
		r := mmap.ReaderAt(26)

		for i := 1; i < 1024; i++ {
			n, err := r.Read(buf)
			tt.AssertIsNotError(t, err)
			tt.AssertEqual(t, 13, n)
			tt.AssertEqual(t, "ABCDEFGHIJKLM", string(buf))

			n, err = r.Read(buf)
			tt.AssertIsNotError(t, err)
			tt.AssertEqual(t, 13, n)
			tt.AssertEqual(t, "NOPQRSTUVWXYZ", string(buf))
		}

		remain, err := io.ReadAll(r)
		tt.AssertIsNotError(t, err)
		tt.AssertEqual(t, make([]byte, mmap.Cap()-len(data)), remain)

		n, err := r.Read(buf)
		tt.AssertEqual(t, io.EOF, err)
		tt.AssertEqual(t, 0, n)
	})
}

package mmap

import (
	"bytes"
	"io"
	"testing"

	"github.com/ImSingee/tt"
)

func TestReadAt(t *testing.T) {
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
}

func TestReadMore(t *testing.T) {
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
}

func TestReadNotExist(t *testing.T) {
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

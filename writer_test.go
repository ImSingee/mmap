package mmap

import (
	"github.com/ImSingee/tt"
	"testing"
)

func TestWrite(t *testing.T) {

}

func TestCopy(t *testing.T) {
	t.Run("copy-closed", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		closeMmap(t, mmap)

		err = mmap.Copy(1, 2, 3)
		tt.AssertEqual(t, ErrIsClosed, err)
	})

	t.Run("no-real-copy", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		err = mmap.Copy(1, 1, 1024) // overflow but ok
		tt.AssertIsNotError(t, err)
	})

	t.Run("copy-overflow", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		// overflow (left)
		err = mmap.Copy(1, 2, LenOfHelloWorld)
		tt.AssertEqual(t, ErrOverflow, err)

		// overflow (right)
		err = mmap.Copy(2, 1, LenOfHelloWorld)
		tt.AssertEqual(t, ErrOverflow, err)
	})

	t.Run("copy-overlap-left", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		err = mmap.Copy(0, 5, LenOfHelloWorld)
		tt.AssertIsNotError(t, err)

		// Hello world!
		// ^    |     ^    |
		//      Hello world!
		// HelloHello world!

		expect := make([]byte, mmap.Cap())
		copy(expect, "HelloHello world!")

		tt.AssertEqual(t, expect, mmap.data)
	})

	t.Run("copy-overlap-right", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		err = mmap.Copy(4, 0, LenOfHelloWorld-4) // "o world!"
		tt.AssertIsNotError(t, err)

		// Hello world!
		// |   ^  |    ^
		// o world!
		// o world!rld!

		expect := make([]byte, mmap.Cap())
		copy(expect, "o world!rld!")

		tt.AssertEqual(t, expect, mmap.data)
	})

	t.Run("copy-more", func(t *testing.T) {
		f, err := newHelloWorldFile()
		tt.AssertIsNotError(t, err)

		mmap, err := New(NewReadWrite(f))
		tt.AssertIsNotError(t, err)
		defer closeMmap(t, mmap)

		err = mmap.Copy(0, 20, LenOfHelloWorld)
		tt.AssertIsNotError(t, err)

		expect := make([]byte, mmap.Cap())
		copy(expect, HelloWorld)
		copy(expect[20:], HelloWorld)

		tt.AssertEqual(t, expect, mmap.data)
	})
}

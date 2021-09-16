//go:build linux || darwin
// +build linux darwin

package mmap

import (
	"golang.org/x/sys/unix"
)

func New(args Opener) (*Mmap, error) {
	if args, ok := args.(shouldClean); ok {
		if err := args.Clean(); err != nil {
			return nil, err
		}
	}

	m := &Mmap{
		args: args,
		grow: DefaultGrowPolicy,
		data: nil,
	}

	return m, m.open(args.InitialSize())
}

type Mmap struct {
	args Opener
	grow Grower

	data   []byte
	closed bool
}

func (m *Mmap) Cap() int {
	return cap(m.data)
}

func (m *Mmap) open(withCap int) (err error) {
	f, err := m.args.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	if size < int64(withCap) {
		err := unix.Ftruncate(int(f.Fd()), int64(withCap))
		if err != nil {
			return err
		}
	}

	m.data, err = unix.Mmap(int(f.Fd()), m.args.Offset(), withCap, m.args.Prot(), m.args.Flags())
	if err == nil {
		m.closed = false
	}
	return
}

func (m *Mmap) IsClosed() bool {
	return m.closed
}

func (m *Mmap) close() (err error) {
	if m.closed {
		return nil
	}

	err = unix.Munmap(m.data)
	m.closed = true

	return
}

func (m *Mmap) reOpen(newCap int) error {
	if err := m.close(); err != nil {
		return err
	}

	return m.open(newCap)
}

func (m *Mmap) EnsureCapacity(size int) error {
	if m.closed {
		return ErrIsClosed
	}

	if capacity := m.Cap(); size > capacity {
		next := m.grow(capacity, size)
		if err := m.reOpen(next); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mmap) Close() error {
	return m.close()
}

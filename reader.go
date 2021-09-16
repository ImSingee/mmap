package mmap

import (
	"bytes"
	"io"
)

var (
	_ io.ReaderAt = (*Mmap)(nil)
	_ io.WriterTo = (*Mmap)(nil)
)

func (m *Mmap) ReadAt(p []byte, off int64) (n int, err error) {
	if m.closed {
		return 0, ErrIsClosed
	}

	if off > int64(len(p)) {
		return 0, io.EOF
	}
	n = copy(p, m.data[off:])
	if n < len(p) {
		err = io.EOF
	}

	return
}

func (m *Mmap) WriteTo(w io.Writer) (n int64, err error) {
	if m.closed {
		return 0, ErrIsClosed
	}

	return io.Copy(w, bytes.NewReader(m.data))
}

func (m *Mmap) Bytes(offset int64, length int) ([]byte, error) {
	if m.closed {
		return nil, ErrIsClosed
	}

	result := make([]byte, length)
	n, err := m.ReadAt(result, offset)
	result = result[:n]
	return result, err
}

package mmap

import "io"

var _ io.WriterAt = (*Mmap)(nil)

func (m *Mmap) WriteAt(p []byte, off int64) (n int, err error) {
	if m.closed {
		return 0, ErrIsClosed
	}

	end := int(off) + len(p)
	if err = m.EnsureCapacity(end); err != nil {
		return 0, err
	}

	return copy(m.data[off:end], p), nil
}

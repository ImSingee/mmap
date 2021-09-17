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

func (m *Mmap) Copy(srcPos, dstPos int64, length int) error {
	if m.closed {
		return ErrIsClosed
	}

	if srcPos == dstPos {
		return nil
	}

	if srcPos+int64(length) > int64(len(m.data)) {
		return ErrOverflow
	}

	if err := m.EnsureCapacity(int(dstPos) + length); err != nil {
		return err
	}

	copy(m.data[dstPos:], m.data[srcPos:srcPos+int64(length)])
	return nil
}

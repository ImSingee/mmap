package mmap

import "fmt"

var ErrIsClosed = fmt.Errorf("mmap is closed")

var ErrOverflow = fmt.Errorf("mmap access out of bound")

package mmap

import (
	"io/ioutil"
	"os"

	"golang.org/x/sys/unix"
)

type Opener interface {
	Open() (*os.File, error)
	Offset() int64
	InitialSize() int
	Prot() int
	Flags() int
}

type shouldClean interface {
	Clean() error
}

type Args struct {
	File       string
	InitLength int
	Readonly   bool
	Private    bool
}

var _ Opener = (*Args)(nil)
var _ shouldClean = (*Args)(nil)

const DefaultInitLength = OneMB

func (a *Args) Clean() error {
	if a.File == "" {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		a.File = f.Name()
		_ = f.Close()

		if !a.Readonly {
			a.InitLength = DefaultInitLength
		}
	}

	if a.InitLength <= 0 {
		n, err := os.Stat(a.File)
		if err != nil {
			return err
		}
		a.InitLength = int(n.Size())
	}

	return nil
}

func (a *Args) Open() (*os.File, error) {
	if a.Readonly {
		return os.Open(a.File)
	} else {
		return os.OpenFile(a.File, os.O_RDWR|os.O_CREATE, 0644)
	}
}

func (a *Args) Offset() int64 {
	return 0
}

func (a *Args) InitialSize() int {
	return a.InitLength
}

func (a *Args) Prot() int {
	if a.Readonly {
		return unix.PROT_READ
	} else {
		return unix.PROT_READ | unix.PROT_WRITE
	}
}

func (a *Args) Flags() int {
	if a.Private {
		return unix.MAP_PRIVATE
	} else {
		return unix.MAP_SHARED
	}
}

func NewReadOnly(file string) *Args {
	return &Args{
		File:       file,
		InitLength: -1,
		Readonly:   true,
		Private:    false,
	}
}

func NewReadWrite(file string) *Args {
	return &Args{
		File:       file,
		InitLength: -1,
		Readonly:   false,
		Private:    false,
	}
}

# mmap for Human

[![Go Reference](https://pkg.go.dev/badge/github.com/ImSingee/mmap.svg)](https://pkg.go.dev/github.com/ImSingee/mmap) [![Test Status](https://github.com/ImSingee/mmap/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/ImSingee/mmap/actions/workflows/test.yml?query=branch%3Amaster) [![codecov](https://codecov.io/gh/ImSingee/mmap/branch/master/graph/badge.svg?token=RWV4ZYS1DH)](https://codecov.io/gh/ImSingee/mmap) [![Go Report Card](https://goreportcard.com/badge/github.com/ImSingee/mmap)](https://goreportcard.com/report/github.com/ImSingee/mmap)

An easy to use mmap wrapper for go. Support read, write and auto grow.

## Installation

```bash
go get -u github.com/ImSingee/mmap
```

## Quick Start

```go
package main

import "github.com/ImSingee/mmap"

func main() {
	f, err := mmap.New(mmap.NewReadWrite(""))
	if err != nil {
		panic("Cannot init mmap: " + err.Error())
	}
	defer f.Close()

	_, err = f.WriteAt([]byte("Hello World!"), 0)
	if err != nil {
		panic(err)
	}
}
```

## License

[MIT License](LICENSE)
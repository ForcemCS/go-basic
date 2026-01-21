package io_test

import (
	"basic/io"
	"testing"
)

// go test -v ./ -run=^TestWriteFile$ -count=1
func TestWriteFile(t *testing.T) {
	io.WriteFile()
}

// go test -v ./ -run=^TestWriteFileWiteBufio$ -count=1
func TestWriteFileWiteBufio(t *testing.T) {
	io.WriteFileWithBufio()
}

// go test -v ./ -run=^TestReadFile$ -count=1
func TestReadFile(t *testing.T) {
	io.ReadFile()
}

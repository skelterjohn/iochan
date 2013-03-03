package iochan

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestBasic(t *testing.T) {
	b := make(Buffer, 10)
	r := bytes.NewBuffer([]byte("this is  a buffer withlongwords sometimes"))
	words, errs := b.ReaderChan(r, " ")
	for s := range words {
		fmt.Printf("%s|", s)
	}
	if err := <-errs; err != nil && err != io.EOF {
		t.Error(err)
	}

	b = make(Buffer, 1024)
	lines, errs := b.FileLineChan("iochan_test.go")
	for s := range lines {
		fmt.Print(s)
	}
	if err := <-errs; err != nil && err != io.EOF {
		t.Error(err)
	}
}

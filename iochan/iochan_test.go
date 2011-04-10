package iochan

import (
	"testing"
	"bytes"
	"fmt"
)

func TestBasic(t *testing.T) {
	b := make(Buffer, 10)
	r := bytes.NewBuffer([]byte("this is  a buffer withlongwords sometimes"))
	for s := range b.ReaderChan(r, " ") {
		fmt.Printf("%s|", s)
	}
	
	b = make(Buffer, 1024)
	for s := range b.FileLineChan("iochan_test.go") {
		fmt.Print(s)
	}
}

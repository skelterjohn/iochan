package iochan

import (
	"testing"
	"bytes"
	"fmt"
)

func TestBasic(t *testing.T) {
	b := make(Buffer, 20)
	ch := b.ReaderChan(bytes.NewBuffer([]byte("this is a buffer withlongwords sometimes")), " ")
	for s := range ch {
		fmt.Println(s)
	}
}

package iochan

import (
	"bytes"
	"io"
	"os"
)

type Buffer []byte

func ReaderChan(r io.Reader, sep string) (cr <-chan string) {
	return make(Buffer, 2048).ReaderChan(r, sep)
}

func (b Buffer) ReaderChan(r io.Reader, sep string) (cr <-chan string) {
	sepb := []byte(sep)
	c := make(chan string)
	go func(cs chan<- string) {
		defer close(cs)
		writeStart := 0
		for {
			if i := bytes.Index(b[:writeStart], sepb); i != -1 {
				msg := string([]byte(b[:i+1]))
				cs <- msg
				copy(b[:writeStart-(i+1)], b[i+1:writeStart])
				writeStart -= i + 1
				continue
			} else if r == nil {
				if writeStart != 0 {
					msg := string([]byte(b[:writeStart]))
					cs <- msg
				}
				break
			} else if writeStart == len(b) {
				msg := string([]byte(b))
				writeStart = 0
				cs <- msg
			}
			if r != nil {
				n, err := r.Read(b[writeStart:])
				if err != nil {
					r = nil
				}
				writeStart += n
			}
		}
	}(c)
	return c
}

func FileLineChan(fpath string) (cr <-chan string) {
	return make(Buffer, 2048).FileLineChan(fpath)
}

func (b Buffer) FileLineChan(fpath string) (cr <-chan string) {
	r, err := os.Open(fpath)
	if err == nil {
		cr = b.ReaderChan(r, "\n")
	} else {
		c := make(chan string)
		close(c)
		cr = c
	}
	return
}

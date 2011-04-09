package iochan

import (
	"io"
	"bytes"
)

type Buffer []byte

func (b Buffer) ReaderChan(r io.Reader, sep string) (cr <-chan string) {
	sepb := []byte(sep)
	c := make(chan string)
	go func(cs chan<- string) {
		writeStart := 0
		for {
			if i := bytes.Index(b[:writeStart], sepb); i != -1 {
				msg := string([]byte(b[:i]))
				cs <- msg
				copy(b[:writeStart-(i+1)], b[i+1:writeStart])
				writeStart -= i+1
				continue
			} else if r == nil {
				msg := string([]byte(b[:writeStart]))
				cs <- msg
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
		close(cs)
	}(c)
	return c
}
package iochan

import (
	"bytes"
	"io"
	"os"
)

type Buffer []byte

func ReaderChan(r io.Reader, sep string) (cr <-chan string, errs <-chan error) {
	return make(Buffer, 2048).ReaderChan(r, sep)
}

func (b Buffer) ReaderChan(r io.Reader, sep string) (cr <-chan string, errs <-chan error) {
	sepb := []byte(sep)
	c := make(chan string)
	e := make(chan error)
	go func(cs chan<- string, errs chan<- error) {
		var rErr error
		writeStart := 0

		defer func() { errs <- rErr; close(errs) }()
		defer close(cs)

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
					rErr = err
				}
				writeStart += n
			}
		}
	}(c, e)
	return c, e
}

func FileLineChan(fpath string) (cr <-chan string, err <-chan error) {
	return make(Buffer, 2048).FileLineChan(fpath)
}

func (b Buffer) FileLineChan(fpath string) (cr <-chan string, errs <-chan error) {
	r, err := os.Open(fpath)
	if err == nil {
		cr, errs = b.ReaderChan(r, "\n")
	} else {
		c := make(chan string)
		close(c)
		cr = c

		e := make(chan error)
		e <- err
		close(e)
	}
	return
}

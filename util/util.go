// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"log"
	"strings"
)

func WriteLengthEncoded(w io.Writer, data []byte) error {
	err := binary.Write(w, binary.BigEndian, int64(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ReadLengthEncoded(r io.Reader) (data []byte, err error) {
	var dataLen int64
	err = binary.Read(r, binary.BigEndian, &dataLen)
	if err != nil {
		return nil, err
	}
	data = make([]byte, dataLen)
	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func JsonPath(object interface{}, path string) interface{} {
	if object == nil {
		return nil
	}

	keys := strings.Split(path, "/")

	if len(keys) == 0 {
		return object
	}

	o := object.(map[string]interface{})
	for i := 0; i < len(keys)-1; i++ {
		o = o[keys[i]].(map[string]interface{})
	}

	return o[keys[len(keys)-1]]
}

func UUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type HashReader struct {
	R io.Reader
	H hash.Hash
}

func (hr *HashReader) Read(p []byte) (n int, err error) {
	n, err = hr.R.Read(p)

	hs := p[0:n]

	hr.H.Write(hs)
	return
}

type AllButTailReader struct {
	ru       io.Reader
	buf      []byte
	tailSize int
	r, w     int
	err      error
}

func NewAllButTailReader(r io.Reader, tailSize int, bufferSize int) *AllButTailReader {
	if bufferSize <= tailSize {
		bufferSize = tailSize + 1
	}
	rv := new(AllButTailReader)
	rv.buf = make([]byte, bufferSize)
	rv.ru = r
	rv.tailSize = tailSize
	return rv
}

func (b *AllButTailReader) Read(p []byte) (n int, err error) {
	slice := b.buf

	n = len(p)
	if n == 0 {
		return 0, b.err
	}

	if b.w <= b.r+b.tailSize {
		if b.err != nil {
			return 0, b.err
		}
		b.fill()
		if b.w <= b.r+b.tailSize {
			return 0, b.err
		}
	}
	if n > b.w-b.r-b.tailSize {
		n = b.w - b.r - b.tailSize
	}
	if n == 0 {
		return 0, b.err
	}
	copy(p[0:n], slice[b.r:])
	b.r += n
	return n, b.err
}

func (b *AllButTailReader) Tail() []byte {
	return b.buf[b.r:b.w]
}

func (b *AllButTailReader) fill() {
	slice := b.buf

	copy(slice, slice[b.r:b.w])
	b.w -= b.r
	b.r = 0

	n, err := b.ru.Read(slice[b.w:])
	b.w += n
	if err != nil {
		b.err = err
	}
}

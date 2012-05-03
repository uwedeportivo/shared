// Copyright (c) 2011 Uwe Hoffmann. All rights reserved.

package steganography

import (
	"codemanic.com/util"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
)

func EncodePNG(w io.Writer, payload []byte, m image.Image) error {
	nrgba := newNRGBAImageLSBReaderWriter(m)

	err := util.WriteLengthEncoded(nrgba, payload)
	if err != nil {
		return err
	}

	return png.Encode(w, nrgba.m)
}

func DecodePNG(rin io.Reader) ([]byte, error) {
	m, err := png.Decode(rin)
	if err != nil {
		return nil, err
	}

	nrgba := newNRGBAImageLSBReaderWriter(m)

	return util.ReadLengthEncoded(nrgba)
}

type nrgbaImageLSBReaderWriter struct {
	m       *image.NRGBA
	x, y, q int
}

func newNRGBAImageLSBReaderWriter(im image.Image) *nrgbaImageLSBReaderWriter {
	rv := new(nrgbaImageLSBReaderWriter)

	rv.x = 0
	rv.y = 0
	rv.q = -1

	b := im.Bounds()

	rv.m = image.NewNRGBA(image.Rect(0, 0, b.Max.X-b.Min.X, b.Max.Y-b.Min.Y))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rv.m.Set(x-b.Min.X, y-b.Min.Y, im.At(x, y))
		}
	}
	return rv
}

func (it *nrgbaImageLSBReaderWriter) reset() {
	it.x = 0
	it.y = 0
	it.q = -1
}

func (it *nrgbaImageLSBReaderWriter) Read(p []byte) (n int, err error) {
	n = 0
	for j, _ := range p {
		var rv byte = 0
		var i uint8
		for i = 0; i < 8; i++ {
			it.q++
			if it.q == 3 {
				it.q = 0
				it.x++
				if it.x == it.m.Rect.Max.X {
					it.x = it.m.Rect.Min.X
					it.y++
					if it.y == it.m.Rect.Max.Y {
						return n, io.EOF
					}
				}
			}

			color := it.m.At(it.x, it.y).(color.NRGBA)
			var colorByte byte
			switch it.q {
			case 0:
				colorByte = color.R
			case 1:
				colorByte = color.G
			case 2:
				colorByte = color.B
			}
			rv = rv | ((colorByte & 1) << i)
		}
		p[j] = rv
		n++
	}
	return n, nil
}

func setLSB(val, bit byte) byte {
	var rv byte
	if bit == 1 {
		rv = val | 1
	} else {
		rv = val & 0xfe
	}
	return rv
}

func (it *nrgbaImageLSBReaderWriter) Write(p []byte) (n int, err error) {
	n = 0
	for _, v := range p {
		var i uint8
		for i = 0; i < 8; i++ {
			it.q++
			if it.q == 3 {
				it.q = 0
				it.x++
				if it.x == it.m.Rect.Max.X {
					it.x = it.m.Rect.Min.X
					it.y++
					if it.y == it.m.Rect.Max.Y {
						return n, io.EOF
					}
				}
			}

			color := it.m.At(it.x, it.y).(color.NRGBA)
			switch it.q {
			case 0:
				color.R = setLSB(color.R, (v>>i)&1)
			case 1:
				color.G = setLSB(color.G, (v>>i)&1)
			case 2:
				color.B = setLSB(color.B, (v>>i)&1)
			}
			it.m.Set(it.x, it.y, color)
		}
		n++
	}
	return n, nil
}

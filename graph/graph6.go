// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"io"
	"io/ioutil"
	"errors"
	"fmt"
	"strings"

	"code.google.com/p/go-bit/bit"
)

const g6Prefix = ">>graph6<<"

func ReadGraph6(r io.Reader) (Graph, error) {
	g6, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(string(g6), g6Prefix) {
		g6 = g6[len(g6Prefix):]
	}

	if len(g6) == 0 {
		return nil, errors.New("no content")
	}

	var n uint64

	if len(g6) > 8 && g6[0] == 126 && g6[1] == 126 {
		n = unr(g6[2:8])
		g6 = g6[8:]
	} else if len(g6) > 4 && g6[0] == 126 {
		n = unr(g6[1:4])
		g6 = g6[4:]
	} else {
		n = uint64(g6[0] - 63)
		g6 = g6[1:]
	}

	if n == 0 {
		return nil, errors.New("no content")
	}

	if n > uint64(bit.MaxInt) {
		return nil, fmt.Errorf("too many vertices %d", n)
	}

	var num int = int(n)

	var g IntAdjacencyListGraph = make([][]Vertex, num)

	ds := make([]int, 6)

	i, j := 0, 1

	ll:
	for _, b := range g6 {
		digits(b, ds)

		for _, v := range ds {
			if v == 1 {
				g.Add(i, j)
				g.Add(j, i)
			}
			if i == num - 2 && j == num - 1 {
				break ll
			}
			i++
			if i == j {
				i, j = 0, j + 1
			}
		}
	}
	return g, nil	    
}

func unr(bs []byte) uint64 {
	var x uint64

	for _, b := range bs {
		x = (x << 6) + uint64(b - 63)
	}
	return x
}

func digits(b byte, ds []int) {
	x := b - 63
	for i := 0; i < 6; i++ {
		ds[5 - i] = int(x & 1)
		x = x >> 1
	}
}
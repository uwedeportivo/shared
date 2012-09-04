// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"os"
	"testing"
)

func TestUnr(t *testing.T) {
	bs := []byte{66, 63, 120}

	x := unr(bs)

	if x != 12345 {
		t.Fatalf("unr failed, expected 12345, got %d", x)
	}

	bs = []byte{63, 90, 90, 90, 90, 90}

	x = unr(bs)

	if x != 460175067 {
		t.Fatalf("unr failed, expected 460175067, got %d", x)
	}
}

func TestDigits(t *testing.T) {
	ds := make([]int, 6)

	digits(81, ds)

	if !SlicesEqual([]int{0,1,0,0,1,0}, ds) {
		t.Fatalf("digits failed")
	}

	digits(99, ds)

	if !SlicesEqual([]int{1,0,0,1,0,0}, ds) {
		t.Fatalf("digits failed")
	}
}

func compareNeighbors(g Graph, u int, e []int) bool {
	expected := make([]Vertex, len(e))

	for i, v := range e {
		expected[i] = IntVertex(v)
	}
	actual := g.Neighbors(IntVertex(u))

	return NeighborsEqual(expected, actual) 
}

func TestGraph6(t *testing.T) {
	file, err := os.Open("testdata/test.g6")
	if err != nil {
		t.Fatalf("error opening test data: %v", err)
	}
	defer file.Close()

	g, err := ReadGraph6(file)

	if err != nil {
		t.Fatalf("failed to read graph6 file: %v", err)
	}

	if g == nil {
		t.Fatalf("failed to read graph6 file")
	}

	if g.NumVertices() != 10 {
		t.Fatalf("number of vertices different from expected value")
	}

	expected := make([][]int, 10)

	expected[0] = []int{1,2,4,6,7,8,9}
	expected[1] = []int{0,2,8,9}
	expected[2] = []int{0,1,5,6}
	expected[3] = []int{8}
	expected[4] = []int{0,7}
	expected[5] = []int{2,6}
	expected[6] = []int{0,2,5,8}
	expected[7] = []int{0,4}
	expected[8] = []int{0,1,3,6,9}
	expected[9] = []int{0,1,8}

	for i := 0; i < 10; i++ {
		if !compareNeighbors(g, i, expected[i]) {
			t.Fatalf("neighbors of vertex %d different from expected value ", i)
		}
	}
}
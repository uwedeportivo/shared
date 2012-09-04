// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"os"
	"testing"
)

func NeighborsEqual(a, b []Vertex) bool {
	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}


func TestReadSnap(t *testing.T) {
	file, err := os.Open("testdata/email-Enron.txt")
	if err != nil {
		t.Fatalf("error opening test data: %v", err)
	}
	defer file.Close()

	g, err := ReadSNAP(file)

	if err != nil {
		t.Fatalf("failed to read snap file: %v", err)
	}

	if g == nil {
		t.Fatalf("failed to read snap file")
	}

	if g.NumVertices() != 36692 {
		t.Fatalf("number of vertices different from expected value")
	}

	expected := []Vertex{IntVertex(1), IntVertex(4), IntVertex(6), IntVertex(878), IntVertex(8552)}
	actual := g.Neighbors(IntVertex(3))

	if !NeighborsEqual(expected, actual) {
		t.Fatalf("neighbors of vertex 3 different from expected value")
	}
}
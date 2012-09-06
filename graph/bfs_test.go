// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"os"
	"testing"
)

func TestBFS(t *testing.T) {
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

	var tree BFSTreeRecorder = make([]int, 10)

	BFS(g, IntVertex(2), []BFSVisitor{tree})

	if !SlicesEqual([]int{2, 2, 0, 8, 0, 2, 2, 0, 0, 0}, tree) {
		t.Fatalf("bfs tree different from expected value")
	}
}
// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

func SlicesEqual(a, b []int) bool {
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


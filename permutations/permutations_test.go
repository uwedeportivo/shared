// Copyright (c) 2013 Uwe Hoffmann. All rights reserved.

package permutations

import (
	"testing"
)

var golden = [][]int{
	{1, 2, 4, 3},
	{1, 4, 2, 3},
	{4, 1, 2, 3},
	{4, 1, 3, 2},
	{1, 4, 3, 2},
	{1, 3, 4, 2},
	{1, 3, 2, 4},
	{3, 1, 2, 4},
	{3, 1, 4, 2},
	{3, 4, 1, 2},
	{4, 3, 1, 2},
	{4, 3, 2, 1},
	{3, 4, 2, 1},
	{3, 2, 4, 1},
	{3, 2, 1, 4},
	{2, 3, 1, 4},
	{2, 3, 4, 1},
	{2, 4, 3, 1},
	{4, 2, 3, 1},
	{4, 2, 1, 3},
	{2, 4, 1, 3},
	{2, 1, 4, 3},
	{2, 1, 3, 4},
}

func slicesEqual(a, b []int) bool {
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

func TestPermutationGenerator(t *testing.T) {
	g, p := NewPermutationGenerator(4)

	for i, _ := range p {
		p[i] = p[i] + 1
	}

	for c, v := range golden {
		i, j, ok := g.Next()

		t.Logf("next: %d, %d, %v\n", i, j, ok)

		if !ok {
			t.Errorf("generating %d permutation failed", c)
		}

		p[i], p[j] = p[j], p[i]

		if !slicesEqual(v, p) {
			t.Errorf("generated permutation %v differs from golden %v", p, v)
		}
	}

	_, _, ok := g.Next()

	if ok {
		t.Errorf("generated more permutations than expected")
	}
}

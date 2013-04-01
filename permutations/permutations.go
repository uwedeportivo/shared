// Copyright (c) 2013 Uwe Hoffmann. All rights reserved.

package permutations

type direction int

const (
	left direction = iota
	right
)

type PermutationGenerator struct {
	p []int
	d []direction
}

func NewPermutationGenerator(size int) (*PermutationGenerator, []int) {
	dp := new(PermutationGenerator)

	dp.p = make([]int, size)
	dp.d = make([]direction, size)

	perm := make([]int, size)

	for i := 0; i < size; i++ {
		dp.p[i] = i
		dp.d[i] = left
		perm[i] = i
	}

	return dp, perm
}

func (dp *PermutationGenerator) swapDirection(i int) {
	if dp.d[i] == left {
		dp.d[i] = right
	} else {
		dp.d[i] = left
	}
}

func (dp *PermutationGenerator) isMobile(i int) bool {
	return (dp.d[i] == left && i > 0 && dp.p[i] > dp.p[i-1]) ||
		(dp.d[i] == right && i < (len(dp.p)-1) && dp.p[i] > dp.p[i+1])
}

func (dp *PermutationGenerator) largestMobile() int {
	max := -1
	maxIndex := -1
	n := len(dp.p)

	for i := 0; i < n; i++ {
		if dp.isMobile(i) && dp.p[i] > max {
			max = dp.p[i]
			maxIndex = i
		}
	}
	return maxIndex
}

func (dp *PermutationGenerator) Next() (ti int, tj int, ok bool) {
	m := dp.largestMobile()

	if m == -1 {
		return -1, -1, false
	}

	mv := dp.p[m]

	if dp.d[m] == left {
		dp.p[m-1], dp.p[m] = dp.p[m], dp.p[m-1]
		dp.d[m-1], dp.d[m] = dp.d[m], dp.d[m-1]
		ti, tj, ok = m-1, m, true
	} else {
		dp.p[m+1], dp.p[m] = dp.p[m], dp.p[m+1]
		dp.d[m+1], dp.d[m] = dp.d[m], dp.d[m+1]
		ti, tj, ok = m, m+1, true
	}

	for i, v := range dp.p {
		if v > mv {
			dp.swapDirection(i)
		}
	}

	return ti, tj, ok
}

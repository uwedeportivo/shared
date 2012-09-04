// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

type sliceQueue struct {
	s []int
	head, tail int
	n int
}

func newQueue(capacity int) *sliceQueue {
	var q sliceQueue
	q.s = make([]int, capacity)
	return &q
}

func (q *sliceQueue) size() int {
	return q.n
}

func (q *sliceQueue) grow() {
	tmp := make([]int, len(q.s) << 1)
	copy(tmp, q.s[q.head:])
	if q.head > 0 {
		copy(tmp[q.head:], q.s[0:q.head])
	}
	q.s = tmp
	q.head = 0
	q.tail = q.n
}

func (q *sliceQueue) enqueue(v int) {
	if q.size() == len(q.s) {
		q.grow()
	}
	q.s[q.tail] = v
	q.tail = (q.tail + 1) % len(q.s)
	q.n++
}

func (q *sliceQueue) dequeue() (int, bool) {
	if q.n == 0 {
		return 0, false
	}
	q.n--
	v := q.s[q.head]
	q.head = (q.head + 1) % len(q.s)
	return v, true
}


// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package graph

import (
	"testing"
)

func TestEmptyQueue(t *testing.T) {
	q := newQueue(10)

	_, ok := q.dequeue()
	if ok {
		t.Fatalf("dequeueing from empty queue succeeded")
	}
}

func TestBasicQueue(t *testing.T) {
	q := newQueue(10)

	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, v := range vs {
		q.enqueue(v)
	}

	for _, v := range vs {
		w, ok := q.dequeue()
		if !ok || v != w {
			t.Fatalf("dequeueing failed %t, expected %d, got %d", ok, v, w)
		}
	}
}

func TestQueueGrowth(t *testing.T) {
	q := newQueue(10)

	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}

	for _, v := range vs {
		q.enqueue(v)
	}

	for _, v := range vs {
		w, ok := q.dequeue()
		if !ok || v != w {
			t.Fatalf("dequeueing failed, expected %d, got %d", v, w)
		}
	}
}

func TestEnqueueDequeue(t *testing.T) {
	q := newQueue(10)

	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}

	for i := 0; i < 10; i++ {
		q.enqueue(vs[i])
	}

	for i := 0; i < 5; i++ {
		w, ok := q.dequeue()
		if !ok || vs[i] != w {
			t.Fatalf("dequeueing failed, expected %d, got %d", vs[i], w)
		}
	}

	for i := 10; i < len(vs); i++ {
		q.enqueue(vs[i])
	}

	for i := 5; i < len(vs); i++ {
		w, ok := q.dequeue()
		if !ok || vs[i] != w {
			t.Fatalf("dequeueing failed, expected %d, got %d", vs[i], w)
		}
	}
}

func TestEnqueueDequeueWrapAround(t *testing.T) {
	q := newQueue(10)

	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	for i := 0; i < 5; i++ {
		q.enqueue(vs[i])
	}

	for i := 0; i < 3; i++ {
		w, ok:= q.dequeue()
		if !ok || vs[i] != w {
			t.Fatalf("dequeueing failed, expected %d, got %d", vs[i], w)
		}
	}

	for i := 5; i < len(vs); i++ {
		q.enqueue(vs[i])
	}

	for i := 3; i < len(vs); i++ {
		w, ok := q.dequeue()
		if !ok || vs[i] != w {
			t.Fatalf("dequeueing failed, expected %d, got %d", vs[i], w)
		}
	}
}

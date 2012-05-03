// Copyright (c) 2012 Uwe Hoffmann. All rights reserved.

package util

import (
	"io"
	"strings"
	"testing"
)

func TestAllButTailReader(t *testing.T) {
	bufferSizes := []int{3, 4, 5, 10, 23, 26, 30}

	for _, bufferSize := range bufferSizes {
		abtr := NewAllButTailReader(strings.NewReader("abcdefghijklmnopqrstuvwxyz"), 3, bufferSize)

		b := make([]byte, 23)
		_, err := io.ReadFull(abtr, b)
		if err != nil {
			t.Fatalf("failed %v", err)
		}

		if string(b) != "abcdefghijklmnopqrstuvw" {
			t.Fatalf("read different from expected %s", string(b))
		}

		if string(abtr.Tail()) != "xyz" {
			t.Fatalf("tail different from expected %s", string(abtr.Tail()))
		}
	}
}

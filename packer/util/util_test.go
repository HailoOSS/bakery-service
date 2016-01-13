package util

import (
	"bufio"
	"fmt"
	"testing"
)

func TestBufferedWriter(t *testing.T) {
	timesCalled := 0

	w := &BufferedWriter{
		WriteFunc: func(p []byte) error {
			timesCalled++
			if len(p) > 10 {
				t.Fatal("Too much written to the writer")
			}

			return nil
		},
	}

	b := bufio.NewWriterSize(w, 10)
	fmt.Fprint(b, "1234567890")

	b.Flush()

	if timesCalled == 0 || timesCalled > 1 {
		t.Fatal("Callback problems")
	}
}

package upload

import (
	"bytes"
	"io"
	"testing"
)

func TestFlushWriterCalled(t *testing.T) {
	wasCalled := false

	fw := NewFlushWriter(1, func(rs io.ReadSeeker) error {
		wasCalled = true
		return nil
	})

	fw.Write([]byte{'1'})

	if !wasCalled {
		t.Fatalf("Flush Writer wasn't called")
	}
}

func TestFlushWriterParts(t *testing.T) {
	fw := NewFlushWriter(5, func(rs io.ReadSeeker) error {
		var buf bytes.Buffer
		io.Copy(&buf, rs)

		if buf.String() != "123456" {
			t.Fatalf("Wrong output back")
		}
		return nil
	})

	fw.Write([]byte{'1', '2', '3', '4'})

	if fw.PartsSent != 0 {
		t.Fatalf("Parts sent should be 0")
	}

	fw.Write([]byte{'5', '6'})

	if fw.PartsSent != 1 {
		t.Fatalf("Parts sent should be 0")
	}
}

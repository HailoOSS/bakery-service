package upload

import (
	"bytes"
	"io"
)

// Flusher:
// 	- Allows for X lines to be written to a buffer
//  - After X lines, flush to somewhere
//  - Allow for a delay, after Y seconds, close.

// WriterFunc is a wrapper function for writers
type WriterFunc func(rs io.ReadSeeker) error

// FlushWriter writes data in bursts
type FlushWriter struct {
	BufferLimit int
	WriterFunc  WriterFunc
	PartsSent   int

	buf bytes.Buffer
}

// NewFlushWriter creates a new FlushWriter
func NewFlushWriter(bufferLimit int, wFn WriterFunc) *FlushWriter {
	return &FlushWriter{
		BufferLimit: bufferLimit,
		WriterFunc:  wFn,
	}
}

// Write data to writer function
func (fw *FlushWriter) Write(p []byte) (int, error) {
	byteLength := len(p)

	if _, err := fw.buf.Write(p); err != nil {
		return 0, err
	}

	if fw.buf.Len() >= fw.BufferLimit {
		if err := fw.WriterFunc(bytes.NewReader(fw.buf.Bytes())); err != nil {
			return byteLength, err
		}

		fw.buf.Reset()
		fw.PartsSent++
	}

	return 0, nil
}

// Len returns the length of the buffer
func (fw *FlushWriter) Len() int {
	return fw.buf.Len()
}

package ui

import (
	"bufio"

	"github.com/hailocab/bakery-service/packer/upload"
)

const (
	// BytesAllowedBeforeFlush marks the amount of bytes allowed to
	// be written before being flushed
	BytesAllowedBeforeFlush = 1024
)

// S3Caller a dummy caller
type S3Caller struct {
	writer *bufio.Writer
}

// NewS3Caller foo
func NewS3Caller(bucket string, path string) *S3Caller {
	w := upload.CallbackWriter{
		WriteFunc: S3PartWrite,
	}

	return &S3Caller{
		writer: bufio.NewWriterSize(w, BytesAllowedBeforeFlush),
	}
}

// Call does something with the message
func (sc *S3Caller) Call(msg *Message) {
	// sc.writer.WriteString()
}

// S3PartWrite writes parts to S3
func S3PartWrite(p []byte) error {
	// TODO: call AWS
	return nil
}

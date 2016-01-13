package ui

import (
	"bufio"
	"bytes"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/packer/util"
)

const (
	// BytesAllowedBeforeFlush marks the amount of bytes allowed to
	// be written before being flushed
	BytesAllowedBeforeFlush = 1024
)

// S3Caller a dummy caller
type S3Caller struct {
	writer *bufio.Writer
	bucket string
	path   string
}

// NewS3Caller foo
func NewS3Caller(bucket string, path string) *S3Caller {
	w := util.CallbackWriter{}

	w.WriteFunc = func(p []byte) error {
		return aws.UploadPart(bucket, path, w.WriteCount, bytes.NewReader(p))
	}

	return &S3Caller{
		writer: bufio.NewWriterSize(w, BytesAllowedBeforeFlush),
	}
}

// Call does something with the message
func (sc *S3Caller) Call(msg *Message) {
	// sc.writer.WriteString()
}

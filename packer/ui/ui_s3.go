package ui

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/hailocab/bakery-service/aws"
	"github.com/hailocab/bakery-service/packer/util"
)

const (
	// BytesAllowedBeforeFlush marks the amount of bytes allowed to
	// be written before being flushed
	BytesAllowedBeforeFlush = 1024
)

var (
	// FlushTimeout is the amount time the buffer will flush
	FlushTimeout = time.Second * 30
)

// S3Caller a dummy caller
type S3Caller struct {
	writer        *bufio.Writer
	bucket        string
	path          string
	activeTimeout bool
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
	sc.writer.WriteString(fmt.Sprintf("%s - %s: %s", time.Now(), msg.Type.String(), msg.Message))

	// Ensure three's a timeout after a msg
	go sc.timeout()
}

func (sc *S3Caller) timeout() {
	if sc.activeTimeout {
		return
	}

	sc.activeTimeout = true

	select {
	case <-time.After(FlushTimeout):
		// Flush buffer and force a write
		sc.writer.Flush()
		sc.activeTimeout = false
	}
}

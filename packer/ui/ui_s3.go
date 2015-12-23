package ui

// S3Caller a dummy caller
type S3Caller struct {
	Uploader Uploader
}

// Call does something with the message
func (sc *S3Caller) Call(msg *Message) {
}

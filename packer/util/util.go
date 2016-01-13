package util

// WriteFunc is a callback type for when a buffer flushes
type WriteFunc func(p []byte) error

// CallbackWriter is a flush buffer to be used with bufio
type CallbackWriter struct {
	WriteCount int64
	WriteFunc  WriteFunc
}

// Write func compat with io.Writer
func (w CallbackWriter) Write(p []byte) (int, error) {
	if err := w.WriteFunc(p); err != nil {
		return 0, err
	}

	w.WriteCount++

	return len(p), nil
}

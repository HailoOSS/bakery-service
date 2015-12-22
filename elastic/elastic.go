// Package elastic is an elasticsearch wrapper
package elastic

// Client wraps the underlying pkg
type Client struct {
	IndexName string
}

// New create new client
func New() *Client {
	return &Client{}
}

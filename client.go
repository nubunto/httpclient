package httpclient

import "net/http"

// Client is the main interface for httpclient
// It is the base of the composable HTTP Client provided.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientFunc is an implementation of Client
// it's used to provide further extension to composable HTTP Clients
type ClientFunc func(*http.Request) (*http.Response, error)

// Do implements the Client interface for ClientFunc
func (cf ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return cf(r)
}

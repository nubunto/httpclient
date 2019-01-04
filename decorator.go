package httpclient

import (
	"encoding/base64"
	"net/http"
	"time"
)

// Decorator is a function that extends a Client with further functionality
type Decorator func(Client) Client

// FaultTolerance applies fault tolerance behaviour to a given client
func FaultTolerance(attempts int, backoff time.Duration) Decorator {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			for i := 0; i < attempts; i++ {
				res, err = c.Do(r)
				if err == nil {
					break
				}
				time.Sleep(backoff * time.Duration(i))
			}
			return
		})
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// BasicAuthorization sets the Authorization header with the basic authorization scheme
func BasicAuthorization(username, password string) Decorator {
	return Authorization("Basic " + basicAuth(username, password))
}

// Authorization adds an Authorization header with the given token
func Authorization(token string) Decorator {
	return Header("Authorization", token)
}

// Header is a generic decorator that adds the provided key/value Header to every request
func Header(key, value string) Decorator {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Add(key, value)
			return c.Do(r)
		})
	}
}

// New gives you a new Client decorated with every Decorator provided
func New(root Client, decorators ...Decorator) Client {
	decorated := root
	for _, decorate := range decorators {
		decorated = decorate(decorated)
	}
	return decorated
}

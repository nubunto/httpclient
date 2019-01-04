package httpclient_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nubunto/httpclient"
)

type testClient struct {
	fail  bool
	calls int
}

func (tc *testClient) Do(r *http.Request) (*http.Response, error) {
	tc.calls++
	if tc.fail {
		return nil, errors.New("failed")
	}
	return &http.Response{}, nil
}

type trackingHandler struct {
	calls int
}

func (th *trackingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.calls++
}

var defaultHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there!"))
})

func newServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	if handler == nil {
		handler = defaultHandler
	}
	return httptest.NewServer(handler)
}

func newRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic("error creating request:" + err.Error())
	}
	return req
}

func TestSimpleClient(t *testing.T) {
	server := newServer(t, nil)
	defer server.Close()

	root := &http.Client{}
	c := httpclient.New(root)
	req := newRequest(t, "GET", server.URL, nil)

	res, err := c.Do(req)
	if err != nil {
		t.Fatal("error should be nil:", err.Error())
	}
	if res.StatusCode != 200 {
		t.Error("should have responded with 200, got", res.StatusCode)
	}
}

func TestFaultTolerance(t *testing.T) {
	server := newServer(t, nil)
	defer server.Close()

	attempts := 2
	root := &testClient{
		fail: true,
	}
	c := httpclient.New(root, httpclient.FaultTolerance(attempts, 500*time.Millisecond))

	req := newRequest(t, "GET", server.URL, nil)
	_, err := c.Do(req)
	if err == nil {
		t.Error("error should be non-nil")
	}
	if root.calls != attempts {
		t.Errorf("client should have been called %d times, actually called %d", attempts, root.calls)
	}
}

package httptest

import (
	"net/http"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

var _ httputil.Requester = (*Client)(nil)

// Request describes a request made to and handled by [Client].
type Request struct {
	// Expectations set the expectations on the request, which are asserted when
	// the request is made.
	Expectations RequestExpectations
	// Response describes the response to the request.
	Response Response
}

// Client asserts that requests made match the expected ones.
type Client struct {
	T        *testing.T
	Requests []Request

	calls int
}

// Resets the calls made to the client.
func (c *Client) Reset() {
	c.calls = 0
}

// Do implements httputil.Requester.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	if c.calls >= len(c.Requests) {
		c.T.Fatalf("Got additional, unexpected request to %s. Did you forget an expectation or to reset the client?", r.URL.String())
	}

	request := c.Requests[c.calls]

	request.Expectations.Assert(c.T, r)

	c.calls++

	return request.Response.Response(c.T, r)
}

// DoCached implements httputil.Requester.
func (c *Client) DoCached(r *http.Request) (*http.Response, error) {
	return c.Do(r)
}

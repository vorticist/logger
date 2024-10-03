package logger

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func Test_LogHttpRequest(t *testing.T) {
	// Define the URL for the request
	u, _ := url.Parse("http://google.com")

	// Define the request body (optional, can be nil for GET requests)
	body := bytes.NewBufferString("request body data")

	// Create a new http.Request object
	req := &http.Request{
		Method: "POST",
		URL:    u,
		Header: http.Header{
			"Content-Type": {"application/json"},
			"User-Agent":   {"Go-Client"},
			"Referer":      {"http://localhost"},
		},
		Body:          io.NopCloser(body),
		ContentLength: int64(body.Len()),
		Host:          "example.com",
	}

	Infof("request: %v", req)
}

package aqi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Setup establishes a test Server that can be used to provide mock responses during testing.
// It returns a pointer to a client, a mux, the server URL and a teardown function that
// must be called when testing is complete.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	c := NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	c.baseURL = url

	return c, mux, server.URL, server.Close
}

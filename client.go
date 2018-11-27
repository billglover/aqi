package aqi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	baseURL   = "https://api.waqi.info/feed"
	userAgent = "github.com/billglover/aqi"
)

// Client holds configuration items for the AQI client and provides methods
// that interact with the AQI API.
type Client struct {
	baseURL   *url.URL
	userAgent string
	client    *http.Client
}

// ClientOptions is a set of options that can be specified when creating an
// AQI client.
type ClientOptions struct {
	BaseURL *url.URL
}

// NewClient returns a new AQI client. If no http.Client is provided then the
// http.DefaultClient is used.
func NewClient(cc *http.Client) *Client {
	if cc == nil {
		cc = http.DefaultClient
	}
	url, _ := url.Parse(baseURL)

	c := &Client{baseURL: url, userAgent: userAgent, client: cc}
	return c
}

// NewClientWithOptions takes ClientOptions, configures and returns a new
// client.
func NewClientWithOptions(cc *http.Client, opts ClientOptions) *Client {
	c := NewClient(cc)
	c.baseURL = opts.BaseURL
	return c
}

// NewRequest creates an HTTP Request. A relative URL should be provided
// without the leading slash. If a non-nil body is provided it will be JSON
// encoded and included in the request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

// Do sends a request and returns the response. An error is returned if the
// request cannot be sent or if the API returns an error. If a response is
// received, the response body is decoded and stored in the value pointed to
// by v.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(err, ctx.Err().Error())
		default:
			return nil, err
		}
	}

	// Anything other than a HTTP 2xx response code is treated as an error.
	if c := resp.StatusCode; c >= 300 {
		err = fmt.Errorf("unexpected response code returned: HTTP %s", resp.Status)
		return resp, err
	}

	// Try and parse the response body.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read body")
	}
	resp.Body.Close()

	// Unmarshal the response body into the provided interface
	if v != nil && len(data) != 0 {
		err = json.Unmarshal(data, v)

		switch err {
		case nil:
		case io.EOF:
			err = nil
		default:
			err = errors.Wrap(err, "unable to parse API response")
		}
	}

	return resp, err
}

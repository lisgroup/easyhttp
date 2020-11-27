package easyhttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client struct
type Client struct {
	options Options
	cli     *http.Client
	req     *http.Request
}

// NewClient func Constructor Client
func NewClient(opts ...Options) *Client {
	var client = Client{}
	if len(opts) > 0 {
		client.options = opts[0]
	}
	return &client
}

// Request Simply encapsulating a HTTP request method
// rawUrl Requested web address
// method Method of request,Only [GET/POST] is supported
// options...
func (c *Client) Request(rawUrl, method string, options ...Options) (response *Response, err error) {
	// 设置 options 并获取 http.Client
	c.setOptions(options...)

	method = strings.ToUpper(method)
	if c.options.BaseURI != "" {
		rawUrl = c.options.BaseURI + rawUrl
	}

	switch method {
	case "GET":
		c.req, err = http.NewRequest("GET", rawUrl, nil)
		if err != nil {
			return
		}
	case "POST":
		c.req, err = http.NewRequest("POST", rawUrl, c.setBody())
		if err != nil {
			return
		}
	default:
		return &Response{}, fmt.Errorf("Only GET and POST methods are supported\n")
	}

	// add http header
	c.setHeader()
	// add query params
	c.setQuery()

	return c.getResponse()
}

// Get request method
func (c *Client) Get(rawUrl string, options ...Options) (*Response, error) {
	return c.Request(rawUrl, "GET", options...)
}

// Post request method
func (c *Client) Post(rawUrl string, options ...Options) (*Response, error) {
	return c.Request(rawUrl, "POST", options...)
}

// setOptions func
func (c *Client) setOptions(options ...Options) {
	if len(options) > 0 {
		newOpts := options[0]
		// Determine whether to replace the public configuration
		if newOpts.BaseURI != "" {
			c.options.BaseURI = newOpts.BaseURI
		}
		if newOpts.Timeout != 0 {
			c.options.Timeout = newOpts.Timeout
		}
		if newOpts.Headers != nil {
			c.options.Headers = newOpts.Headers
		}
		if newOpts.BodyMaps != nil {
			c.options.BodyMaps = newOpts.BodyMaps
		}
		if newOpts.Query != nil {
			c.options.Query = newOpts.Query
		}
		if newOpts.JSON != nil {
			c.options.JSON = newOpts.JSON
		}
	}
	// default 20 second
	if c.options.Timeout <= 0 {
		c.options.Timeout = 20
	}
	// skip verify SSL certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c.cli = &http.Client{
		Timeout:   time.Duration(c.options.Timeout*1000) * time.Millisecond,
		Transport: tr,
	}
}

// setHeader add http header
func (c *Client) setHeader() {
	for key, val := range c.options.Headers {
		c.req.Header.Set(key, fmt.Sprintf("%v", val))
	}
}

// setQuery
func (c *Client) setQuery() {
	switch c.options.Query.(type) {
	case string:
		str := c.options.Query.(string)
		c.req.URL.RawQuery = str
	case map[string]interface{}:
		q := c.req.URL.Query()
		for k, v := range c.options.Query.(map[string]interface{}) {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		c.req.URL.RawQuery = q.Encode()
	}
}

// setBody
func (c *Client) setBody() io.Reader {
	// Form => application/x-www-form-urlencoded
	if c.options.BodyMaps != nil {
		values := url.Values{}
		for k, v := range c.options.BodyMaps {
			if vv, ok := v.(string); ok {
				values.Set(k, vv)
			}
			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					values.Add(k, vvv)
				}
			}
		}
		return strings.NewReader(values.Encode())
	}

	// json => application/json
	if c.options.JSON != nil {
		b, err := json.Marshal(c.options.JSON)
		if err == nil {
			return bytes.NewReader(b)
		}
	}
	return nil
}

// getResponse Processing returned results
func (c *Client) getResponse() (response *Response, err error) {
	startTime := time.Now() // start time
	resp, err := c.cli.Do(c.req)
	cost := time.Since(startTime) // cost time
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	//     err = fmt.Errorf("Get content failed status code is %d ", resp.StatusCode)
	//     return
	// }
	var result string
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		// The buf data read by the cycle is accumulated and stored in the result
		result += string(buf[:n])
	}
	// return Response struct
	return &Response{
		Content:  result,
		HttpCode: resp.StatusCode,
		Status:   resp.Status,
		Cost:     float32(cost),
	}, nil
}

package easyhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client 结构体
type Client struct {
	options Options
}

// Options 结构体
type Options struct {
	BaseURI  string
	Timeout  float32
	Headers  map[string]interface{}
	BodyMaps map[string]interface{}
}

// Client 构造函数
func NewClient(opts ...Options) *Client {
	var client = Client{}
	if len(opts) > 0 {
		client.options = opts[0]
	}
	return &client
}

// Request 简单封装一个Http请求方法
// rawUrl 请求的接口地址
// method 请求的方法，GET/POST
// options 可选参数
func (c *Client) Request(rawUrl, method string, options ...Options) (result string, err error) {
	if len(options) > 0 {
		newOpts := options[0]
		// 判断是否替换公共配置
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
	}
	// 默认超时 20 秒
	if c.options.Timeout <= 0 {
		c.options.Timeout = 20
	}

	client := &http.Client{
		Timeout: time.Duration(c.options.Timeout*1000) * time.Millisecond,
	}
	// 请求的 body 内容
	// 判断headers头信息是否存在 content-type:application/json
	var isJson bool
	for key, value := range c.options.Headers {
		val := fmt.Sprintf("%v", value)
		if strings.ToLower(key) == "content-type" && strings.ToLower(val) == "application/json" {
			isJson = true
			break
		}
	}

	var sendBody io.Reader
	if isJson {
		res, errJson := json.Marshal(c.options.BodyMaps)
		if errJson != nil {
			err = errJson
			return
		}
		sendBody = bytes.NewBuffer(res)
	} else {
		data := url.Values{}
		for key, value := range c.options.BodyMaps {
			// interface 转 string
			data.Set(key, fmt.Sprintf("%v", value))
		}
		sendBody = strings.NewReader(data.Encode())
	}

	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		method = "GET"
	}
	if c.options.BaseURI != "" {
		rawUrl = c.options.BaseURI + rawUrl
	}
	// 提交请求
	request, err := http.NewRequest(method, rawUrl, sendBody) // URL-encoded|JSON payload
	if err != nil {
		return
	}
	// 增加header头信息
	for key, val := range c.options.Headers {
		request.Header.Set(key, fmt.Sprintf("%v", val))
	}
	// 处理返回结果
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Get content failed status code is %d ", response.StatusCode)
	}
	buf := make([]byte, 4096)
	for {
		n, err2 := response.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		// 累加循环读取的 buf 数据，存入 result 中
		result += string(buf[:n])
	}
	return
}

// Get 请求方式
func (c *Client) Get(rawUrl string, options ...Options) (result string, err error) {
	return c.Request(rawUrl, "GET", options...)
}

// Post 请求方式
func (c *Client) Post(rawUrl string, options ...Options) (result string, err error) {
	return c.Request(rawUrl, "POST", options...)
}

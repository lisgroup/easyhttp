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
	BaseURI  string
	Timeout  float32
	Headers  map[string]interface{}
	BodyMaps map[string]interface{}
}

// Request 简单封装一个Http请求方法
// rawUrl 请求的接口地址
// method 请求的方法，GET/POST
// opt 可选参数
func (c *Client) Request(rawUrl, method string, cli ...Client) (result string, err error) {
	if len(cli) > 0 {
		c.BaseURI = cli[0].BaseURI
		c.Timeout = cli[0].Timeout
		c.Headers = cli[0].Headers
		c.BodyMaps = cli[0].BodyMaps
	}
	// 默认超时 20 秒
	if c.Timeout <= 0 {
		c.Timeout = 20
	}

	client := &http.Client{
		Timeout: time.Duration(c.Timeout*1000) * time.Millisecond,
	}
	// 请求的 body 内容
	// 判断headers头信息是否存在 content-type:application/json
	var isJson bool
	for key, value := range c.Headers {
		val := fmt.Sprintf("%v", value)
		if strings.ToLower(key) == "content-type" && strings.ToLower(val) == "application/json" {
			isJson = true
		}
	}

	var sendBody io.Reader
	if isJson {
		res, errJson := json.Marshal(c.BodyMaps)
		if errJson != nil {
			err = errJson
			return
		}
		sendBody = bytes.NewBuffer(res)
	} else {
		data := url.Values{}
		for key, value := range c.BodyMaps {
			// interface 转 string
			data.Set(key, fmt.Sprintf("%v", value))
		}
		sendBody = strings.NewReader(data.Encode())
	}

	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		method = "GET"
	}
	if c.BaseURI != "" {
		rawUrl = c.BaseURI + rawUrl
	}
	// 提交请求
	request, err1 := http.NewRequest(method, rawUrl, sendBody) // URL-encoded|JSON payload
	if err1 != nil {
		err = err1
		return
	}
	// 增加header头信息
	for key, val := range c.Headers {
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
func (c *Client) Get(rawUrl string, cli ...Client) (result string, err error) {
	return c.Request(rawUrl, "GET", cli...)
}

// Post 请求方式
func (c *Client) Post(rawUrl string, cli ...Client) (result string, err error) {
	return c.Request(rawUrl, "POST", cli...)
}

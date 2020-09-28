package easyhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 简单封装一个Http请求方法
// rawUrl 请求的接口地址
// method 请求的方法，GET/POST
// bodyMap 请求的 body 内容
// header 请求的头信息
// timeout 超时时间
func Request(rawUrl, method string, bodyMaps, headers map[string]string, timeout time.Duration) (result string, err error) {
	if timeout <= 0 {
		timeout = 5
	}
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	// 请求的 body 内容
	data := url.Values{}
	for key, value := range bodyMaps {
		data.Set(key, value)
	}
	// 提交请求
	request, err1 := http.NewRequest(method, rawUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err1 != nil {
		err = err1
		return
	}
	// 增加header头信息
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	// 处理返回结果
	response, _ := client.Do(request)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get content failed status code is %d ", response.StatusCode)
	}
	res, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		err = err2
		return
	}
	return string(res), nil
}

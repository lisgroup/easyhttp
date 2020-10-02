package easyhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 简单封装一个Http请求方法
// rawUrl 请求的接口地址
// method 请求的方法，GET/POST
// bodyMap 请求的 body 内容
// header 请求的头信息
// timeout 超时时间
func Request(rawUrl, method string, bodyMaps map[string]interface{}, headers map[string]string, timeout time.Duration) (result string, err error) {
	if timeout <= 0 {
		timeout = 5
	}
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	// 请求的 body 内容
	// 判断headers头信息是否存在content-type:application/json
	var isJson bool
	for key, value := range headers {
		if strings.ToLower(key) == "content-type" && strings.ToLower(value) == "application/json" {
			isJson = true
		}
	}

	var sendBody io.Reader
	if isJson {
		res, errJson := json.Marshal(bodyMaps)
		if errJson != nil {
			err = errJson
			return
		}
		sendBody = bytes.NewBuffer(res)
	} else {
		data := url.Values{}
		for key, value := range bodyMaps {
			// interface 转 string
			str, errInter := interfaceToString(value)
			if errInter != nil {
				err = errInter
				return
			}
			data.Set(key, str)
		}
		sendBody = strings.NewReader(data.Encode())
	}

	// 提交请求
	request, err1 := http.NewRequest(method, rawUrl, sendBody) // URL-encoded|JSON payload
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

func interfaceToString(unknown interface{}) (str string, err error) {
	// interface 转 string
	switch unknown.(type) {
	case int:
		str = strconv.Itoa(unknown.(int))
	case string:
		str = unknown.(string)
	case byte:
		str = string(unknown.([]byte))
	default:
		err = fmt.Errorf("type: %T not supported; Only support int/string/[]byte", unknown)
	}
	return str, err
}

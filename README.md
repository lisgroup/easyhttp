# http

#### 介绍
简易封装 http 请求包

#### 使用说明

1. 获取安装包

`go get -u github.com/lisgroup/easyhttp`

2.  直接使用
```go
package main

import (
    "fmt"

    "github.com/lisgroup/easyhttp"
)

func main() {
	header := map[string]interface{}{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
	}
	bodyMaps := map[string]interface{}{
		"key1": "1111",
		"key2": "222",
	}
	cli := easyhttp.Client{BaseURI: "", Timeout: 4.00, Headers: header, BodyMaps: bodyMaps}
	
	result, err := cli.Request("http://localhost/", "get")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

```

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
	body := map[string]string{}
    resp, _ := easyhttp.Request("http://www.example.com/", "GET", body, body, 5)
    fmt.Println(resp)
}
```

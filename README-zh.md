# easyhttp

#### 介绍
简易封装 http 请求包

#### 使用说明

1. 获取安装包

`go get -u github.com/lisgroup/easyhttp`

2.  直接使用
##### 2.1 Get 方式
```go
package main

import (
    "fmt"
    "log"

    "github.com/lisgroup/easyhttp"
)

func main() {
    cli := easyhttp.NewClient()
    
    result, err := cli.Get("http://localhost/")
    if err != nil {
        log.Println(err)
    }
    fmt.Println(result.Content)
}

```

##### 2.2 Post 方式
```go
func main() {
    cli := easyhttp.NewClient()

    result, err := cli.Post("http://localhost/")
    if err != nil {
        log.Println(err)
    }
    fmt.Println(result.Content)
}
```

##### 2.3 带 Header 和数据的请求
```go
func main() {
    header := map[string]interface{}{
        "Content-Type": "application/json",
        "User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36",
    }
    bodyMaps := map[string]interface{}{
        "key1": "1111",
        "key2": "222",
    }
    options := easyhttp.Options{BodyMaps: bodyMaps, Headers: header, Timeout: 4}
    cli := easyhttp.NewClient(options)

    result, err := cli.Post("http://localhost/")
    if err != nil {
        log.Println(err)
    }
    fmt.Println(result.Content)
}
```

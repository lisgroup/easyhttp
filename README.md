# easyhttp

#### Introduce
Simple encapsulation of HTTP request package

#### Usage

1. Download and install

`go get -u github.com/lisgroup/easyhttp`

2.  Demo
##### 2.1 Get method
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

##### 2.2 Post method
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

##### 2.3 add Header and Request BodyMaps
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

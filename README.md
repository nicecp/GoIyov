# GoIyov

#### 介绍
golang 代理


#### 安装教程
```bash
 go get -u github/nicecp/GoIyov
```

#### 支持特性
* HTTP/HTTPS代理
* TLS/SSL解密
* MITM(中间人攻击)
* 自定义DNS
* Certificate缓存
* Statistic统计(开发中)

![软件结构图](docs/GoIyov.jpg)

#### 证书安装
```
$ go run main.go [-cert]
    -cert 安装服务端证书，支持MacOS(需root权限)/Windows系统
```

#### 使用说明
##### 代理
```go
package main
import (
	"github.com/nicecp/GoIyov"
	"net/http"
	"time"
)

func main() {
	proxy := GoIyov.New()
	server := &http.Server{
		Addr:         ":8888",
		Handler:	  http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			proxy.ServerHandler(rw, req)
		}),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
```

##### 自定义DNS
```go
proxy.AddDnsRecords(map[string]string{"localhost-x":"127.0.0.1"})
```
##### MITM(中间人攻击)
```go
package main

import (
	"fmt"
	"github.com/nicecp/GoIyov"
	"github.com/nicecp/GoIyov/entity"
	"net/http"
	"time"
)

type Handler struct {
	GoIyov.Delegate
}

func (handler *Handler) BeforeRequest(entity *entity.Entity) {
	fmt.Printf("%+v",entity.GetRequestBody())
}
func (handler *Handler) BeforeResponse(entity *entity.Entity, err error) {
	fmt.Printf("%+v",entity.GetResponseBody())
}
func (handler *Handler) ErrorLog(err error) {}

func main() {
	proxy := GoIyov.NewWithDelegate(&Handler{})
	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			proxy.ServerHandler(rw, req)
		}),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
```
##### 移动端如何调试
请确保手机及电脑连接同一个局域网，并将移动设备HTTP代理设置为本机`IP:PORT`
1. 浏览器打开 `goiyov.io/ssl`
2. 安装并信任证书
> 建议使用`MITM`特性，以便查看明文内容

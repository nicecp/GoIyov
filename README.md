# GoIyov

#### 介绍
golang 代理


#### 安装教程
```bash
 go get github/nicecp/IyovGo
```

#### 支持特性
* HTTP/HTTPS代理
* TLS/SSL解密
* MITM(中间人攻击)
* 自定义DNS
* Certiface缓存
* Statistic统计(开发中)

#### 使用说明
![软件结构图](docs/IyovGo.jpg)
> ***双击 `cert/caRoot.crt`根证书文件，并信任该证书***

```go
import (
	"IyovGo"
	"net/http"
	"time"
)

func main() {
	proxy := new(IyovGo.Proxy)
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

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)

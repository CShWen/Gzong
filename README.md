# gzong 

[![Sourcegraph](https://sourcegraph.com/github.com/cshwen/gzong/-/badge.svg)](https://sourcegraph.com/github.com/cshwen/gzong)
[![GoDoc](https://godoc.org/github.com/cshwen/gzong?status.svg)](https://godoc.org/github.com/cshwen/gzong)
[![Go Report Card](https://goreportcard.com/badge/github.com/cshwen/gzong)](https://goreportcard.com/report/github.com/cshwen/gzong)
[![Build Status](http://img.shields.io/travis/cshwen/gzong.svg)](https://travis-ci.org/cshwen/gzong)
[![Codecov](https://img.shields.io/codecov/c/github/cshwen/gzong.svg)](https://codecov.io/gh/cshwen/gzong)
[![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/cshwen/gzong/master/LICENSE)


## Installation
一个以Go为基础的简单web框架，支持中间件插入，已有3个简单的中间件可集成（basicAuth、logging、session），部分完成单元测试。

demo中简单实现了支持增删查改的用户中心(docker+MongoDB)仅供参考。



## Usage

试着写一个hello world，可参照gzong/example/easy.go

```go
package main

import (
	"fmt"
	"github.com/cshwen/gzong"
	"net/http"
)

func main() {
	gz := gzong.New()
	gz.GET("/test", testFunc)
	gz.Run(":8080")
}

func testFunc(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello gzong.")
}
```

然后写完记得先编译下（go build）再启动服务(go run)

```sh
# 编译
$ go build gzong
```

```sh
# 启动服务
$ go run gzong/example/easy.go
```

```sh
# 测试服务
$ curl localhost:8080/test
```

Ps: 还有个已实现好的比较复杂示例gzong/example/complex.go可参照上述流程实操试验下。



## TODO

1. 基础web服务，可启动可正常访问 ✔️

2. 服务框架化 ✔️

3. 路由支持 ✔️

4. 中间件支持 ✔️

5. logging支持 ✔️

6. basicAuth支持 ✔️

7. session支持 ✔️

8. demo实现 ✔️

9. MongoDB环境(docker) ✔️

10. MongoDB CRUD实现(用户中心) ✔️

11. 测试 ✔️

12. 文档 ✔️

13. 汇总 ✔️

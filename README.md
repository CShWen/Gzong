# gzong 

[![Sourcegraph](https://sourcegraph.com/github.com/cshwen/gzong/-/badge.svg)](https://sourcegraph.com/github.com/cshwen/gzong)
[![GoDoc](https://godoc.org/github.com/cshwen/gzong?status.svg)](https://godoc.org/github.com/cshwen/gzong)
[![Go Report Card](https://goreportcard.com/badge/github.com/cshwen/gzong)](https://goreportcard.com/report/github.com/cshwen/gzong)
[![Build Status](http://img.shields.io/travis/cshwen/gzong.svg)](https://travis-ci.org/cshwen/gzong)
[![Codecov](https://img.shields.io/codecov/c/github/cshwen/gzong.svg)](https://codecov.io/gh/cshwen/gzong)
[![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/cshwen/gzong/master/LICENSE)

一个以Go为基础的简单web框架，简单实现了支持增删查改的用户中心demo（MongoDB）。



## Examples

```sh
# 编译
$ go build gzong
```

```sh
# 启动服务
$ go run gzong/example/hello.go
```

```sh
# 测试服务
$ curl localhost:8080/test
```



## 用户中心demo

简单实现了支持增删查改的用户中心demo，数据库实例为docker上的MongoDB镜像。

```sh
# 运行环境 docker+MongoDB镜像
docker run -p 27017:27017 -v /tmp/db:/data/db -d mongo
# 启动服务
$ go build gzong
$ go run gzong/demo/demo.go
# 创建新用户
$ curl localhost:8080/addUser -X POST -d 'nick=tnick&password=tpwd&name=tname&email=tt@gmail.com&phone=13712345678'
# 修改用户信息
$ curl localhost:8080/updateUser -X POST -d 'nick=tnick&password=tpwd&name=tname&email=tt@gmail.com&phone=13712345678'
# 查询用户
$ curl localhost:8080/queryUser -X POST -d 'nick=tnick&name=tname'
# 删除用户
$ curl localhost:8080/delUser -X POST -d 'nick=tnick&name=tname'

```



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

# gzong

gzong，一个以Go为基础的简单web框架



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

11. 测试

12. 文档

13. 汇总

    

##### 建议/待优化点 

1. router的/test改为单元测试
  https://github.com/cshwen/gzong/blob/master/gzong.go#L18
  https://github.com/cshwen/gzong/blob/master/gzong.go#L22
  例如：
  https://github.com/labstack/echo/blob/master/bind_test.go

2. 尽量避免 relative import
  https://github.com/cshwen/gzong/blob/master/example/hello.go#L6

3. 代码格式化 go fmt 

4. 文档， report card， ci, 测试覆盖率
https://goreportcard.com/report/github.com/cshwen/gzong
参考：
https://github.com/labstack/echo#guide

5. 补充文档 https://godoc.org/github.com/cshwen/gzong

6. GitHub改名cshwen，项目改名gzong
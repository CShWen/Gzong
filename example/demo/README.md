# 用户中心demo

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
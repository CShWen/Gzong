package main

import (
	"github.com/cshwen/gzong"
	"github.com/cshwen/gzong/demo/uc"
)

func main() {
	uic := new(uc.UserCenter)
	uic.ConnectToDB("mongodb://localhost:27017", "gz", "user")

	gz := gzong.New()
	gz.POST("/addUser", uic.AddUserFunc)
	gz.POST("/updateUser", uic.UpdateUserFunc)
	gz.POST("/queryUser", uic.QueryUserFunc)
	gz.POST("/delUser", uic.DelUserFunc)

	gz.Run(":8080")
}

package uc

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gzong/demo/mongo"
	"net/http"
	"time"
)

// UserCenter 用户中心的配置和MongoDB核心
type UserCenter struct {
	dbc            *mongo.GzMongo
	dbName         string
	collectionName string
}

// BaseUser 用户标元数据
type BaseUser struct {
	_id      string `bson:",omitempty"`
	Nick     string `bson:",omitempty"`
	Password string `bson:",omitempty"`
	Name     string `bson:",omitempty"`
	Email    string `bson:",omitempty"`
	Phone    string `bson:",omitempty"`
}

// ConnectToDB 连接指定MongoDB
func (uc *UserCenter) ConnectToDB(uri, dbName, collectionName string) {
	uc.dbc = new(mongo.GzMongo)
	uc.dbc.ConnectDB(uri)
	uc.dbName = dbName
	uc.collectionName = collectionName
}

//  md5加密内容
func md5Encryption(content string) string {
	if len(content) == 0 {
		return ""
	}

	data := []byte(content)
	has := md5.Sum(data)
	md5content := fmt.Sprintf("%x", has)
	return md5content
}

// AddUserFunc 添加新用户
func (uc UserCenter) AddUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	newUser := &BaseUser{
		Nick:     req.PostForm.Get("nick"),
		Password: md5Encryption(req.PostForm.Get("password")),
		Name:     req.PostForm.Get("name"),
		Email:    req.PostForm.Get("email"),
		Phone:    req.PostForm.Get("phone"),
	}

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, newUser)
	if err == nil {
		id := res.InsertedID
		fmt.Fprintf(w, "add user success. InsertedID:%s \n", id)
	} else {
		fmt.Fprintf(w, "add user failed. %s \n", err)
	}
}

// DelUserFunc 可根据nick或name删除用户
func (uc UserCenter) DelUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	delUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
		Name: req.PostForm.Get("name"),
	}
	delUserBson, _ := bson.Marshal(delUser)

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, delUserBson)

	if err == nil {
		fmt.Fprintf(w, "delete user success. DeletedCount:%d \n", res.DeletedCount)

	} else {
		fmt.Fprintf(w, "delete user failed. %s \n", err)
	}
}

// QueryUserFunc 可根据nick或name查询用户信息
func (uc UserCenter) QueryUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var userResult BaseUser
	filterUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
		Name: req.PostForm.Get("name"),
	}
	filterUserBson, err := bson.Marshal(filterUser)
	if err != nil {
		fmt.Fprintf(w, "query user failed. %s \n ", err)
	}

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, filterUserBson).Decode(&userResult)
	if err == nil {
		jsonBytes, _ := json.Marshal(userResult)
		fmt.Fprintln(w, string(jsonBytes))
	} else {
		fmt.Fprintf(w, "query user failed. %s \n", err)
	}
}

// UpdateUserFunc 根据nick去修改对应用户信息
func (uc UserCenter) UpdateUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	filterUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
	}
	updateUser := &BaseUser{
		Password: md5Encryption(req.PostForm.Get("password")),
		Name:     req.PostForm.Get("name"),
		Email:    req.PostForm.Get("email"),
		Phone:    req.PostForm.Get("phone"),
	}
	filterUserBson, _ := bson.Marshal(filterUser)

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, filterUserBson, bson.M{"$set": updateUser})

	if err == nil {
		fmt.Fprintf(w, "update user success. %s \n", res)
	} else {
		fmt.Fprintf(w, "update user failed. %s \n", err)
	}
}

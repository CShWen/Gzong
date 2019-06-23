package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gzong"
	"gzong/mongo"
	"net/http"
	"time"
)

func main() {
	uic := new(UserCenter)
	uic.ConnectToDB("mongodb://localhost:27017", "gz", "user")

	gz := gzong.New()
	gz.POST("/addUser", uic.addUserFunc)
	gz.POST("/updateUser", uic.updateUserFunc)
	gz.POST("/queryUser", uic.queryUserFunc)
	gz.POST("/delUser", uic.delUserFunc)

	gz.Run(":8080")
}

type UserCenter struct {
	dbc            *mongo.GzMongo
	dbName         string
	collectionName string
}

type BaseUser struct {
	_id      string `bson:",omitempty"`
	Nick     string `bson:",omitempty"`
	Password string `bson:",omitempty"`
	Name     string `bson:",omitempty"`
	Email    string `bson:",omitempty"`
	Phone    string `bson:",omitempty"`
}

func (uc *UserCenter) ConnectToDB(uri, dbName, collectionName string) {
	uc.dbc = new(mongo.GzMongo)
	uc.dbc.ConnectDB(uri)
	uc.dbName = dbName
	uc.collectionName = collectionName
}

func (uc UserCenter) ucLoginFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "mongo test.\n")
}
func (uc UserCenter) ucLogoutFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "mongo test.\n")
}

//  md5加密内容
func encryption(content string) string {
	if len(content) == 0 {
		return ""
	}

	data := []byte(content)
	has := md5.Sum(data)
	md5content := fmt.Sprintf("%x", has)
	return md5content
}

func (uc UserCenter) addUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	newUser := &BaseUser{
		Nick:     req.PostForm.Get("nick"),
		Password: encryption(req.PostForm.Get("password")),
		Name:     req.PostForm.Get("name"),
		Email:    req.PostForm.Get("email"),
		Phone:    req.PostForm.Get("phone"),
	}

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, newUser)
	if err == nil {
		id := res.InsertedID
		fmt.Fprintf(w, "add user success. ", id)
	} else {
		fmt.Fprintf(w, "add user failed. ", err)
	}
}

func (uc UserCenter) delUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	delUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
		Name: req.PostForm.Get("name"),
	}
	delUserBson, _ := bson.Marshal(delUser)

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.DeleteOne(ctx, delUserBson)

	if err == nil {
		fmt.Fprintf(w, "delete user success. ", res.DeletedCount)

	} else {
		fmt.Fprintf(w, "delete user failed. ", err)
	}
}

func (uc UserCenter) queryUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var userResult BaseUser
	filterUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
		Name: req.PostForm.Get("name"),
	}
	filterUserBson, err := bson.Marshal(filterUser)
	if err != nil {
		fmt.Fprintf(w, "query user failed. ", err)
	}

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, filterUserBson).Decode(&userResult)
	if err == nil {
		jsonBytes, _ := json.Marshal(userResult)
		fmt.Fprintf(w, string(jsonBytes))
	} else {
		fmt.Fprintf(w, "query user failed. ", err)
	}
}

func (uc UserCenter) updateUserFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	filterUser := &BaseUser{
		Nick: req.PostForm.Get("nick"),
	}
	updateUser := &BaseUser{
		Password: encryption(req.PostForm.Get("password")),
		Name:     req.PostForm.Get("name"),
		Email:    req.PostForm.Get("email"),
		Phone:    req.PostForm.Get("phone"),
	}
	filterUserBson, _ := bson.Marshal(filterUser)

	collection := uc.dbc.GetCollection(uc.dbName, uc.collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.UpdateOne(ctx, filterUserBson, bson.M{"$set": updateUser})

	if err == nil {
		fmt.Fprintf(w, "update user success. ", res)
	} else {
		fmt.Fprintf(w, "update user failed. ", err)
	}
}

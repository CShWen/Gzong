package uc

import (
	"gzong/mongo"
)

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

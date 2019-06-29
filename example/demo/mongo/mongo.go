package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// GzMongo MongoDB关键部件
type GzMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
	err        error
}

// ConnectDB 连接到指定MongoDB服务
func (mc *GzMongo) ConnectDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mc.client, mc.err = mongo.NewClient(options.Client().ApplyURI(uri))
	if mc.err != nil {
		log.Fatal(mc.err)
	}
	mc.err = mc.client.Connect(ctx)
	if mc.err != nil {
		log.Fatal(mc.err)
	}

	mc.err = mc.client.Ping(context.TODO(), nil)
	if mc.err != nil {
		log.Fatal(mc.err)
	}
	fmt.Println("Connected to DB.")
}

// GetCollection 连接后返回指定库表的表对象
func (mc *GzMongo) GetCollection(dbName, collectionName string) *mongo.Collection {
	if mc.collection == nil {
		mc.collection = mc.client.Database(dbName).Collection(collectionName)
	}
	return mc.collection
}

package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type GzMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
	err        error
}

func (mc *GzMongo) ConnectDB(uri string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

func (mc *GzMongo) GetCollection(dbName, collectionName string) *mongo.Collection {
	if mc.collection == nil {
		mc.collection = mc.client.Database(dbName).Collection(collectionName)
	}
	return mc.collection
}

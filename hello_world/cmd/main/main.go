package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const helloWorldField = "helloWorldField"

func main() {
	client := prepareClient()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testdb").Collection("testcollection")
	documents := bson.D{{Key: helloWorldField, Value: "Hello, World!"}}
	res, err := collection.InsertOne(context.Background(), documents)
	if err != nil {
		panic(fmt.Errorf("can't insert document: %w", err))
	}
	fmt.Printf("InsertedId: %v, Acknowledged: %v\n", res.InsertedID, res.Acknowledged)

	var result bson.M
	filter := bson.M{helloWorldField: "Hello, World!"}
	if err = collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		panic(fmt.Errorf("can't find document: %w", err))
	}
	fmt.Println(result[helloWorldField])
}

func prepareClient() *mongo.Client {
	connectUri := os.Getenv("MONGO_URI")
	if connectUri == "" {
		panic("MONGO_URI env must be specified")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(connectUri))
	if err != nil {
		panic(fmt.Errorf("can't connect to MongoDB: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(fmt.Errorf("can't ping MongoDB: %w", err))
	}

	return client
}

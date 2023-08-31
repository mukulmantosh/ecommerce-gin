package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func DBSet() *mongo.Client {
	connectionString := "mongodb://localhost:27017"
	clientOption := options.Client().ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println("MongoDB Connection Success!")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {

}

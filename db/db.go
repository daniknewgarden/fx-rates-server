package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMongoClient returns a MongoDB client
func GetMongoClient() *mongo.Client {
	credential := options.Credential{
   		Username: "root",
   		Password: "mySecureDbPassword1!",
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
    	log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
    	log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
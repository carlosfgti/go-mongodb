package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:pass@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	InsertUser(client)
	GetAllUsers(client)
}

func InsertUser(client *mongo.Client) {
	collection := client.Database("golang").Collection("users")
	user := bson.D{
		{"name", "John Doe"},
		{"age", 29},
		{"email", "john@go.dev"},
	}
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted user with ID:", result.InsertedID)
}

func GetAllUsers(client *mongo.Client) {
	collection := client.Database("golang").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			panic(err)
		}
		fmt.Printf("User: %+v\n", user["name"])
	}
}

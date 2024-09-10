package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	id := InsertUser(client)
	GetAllUsers(client)
	GetUserById(client, id)
}

func InsertUser(client *mongo.Client) string {
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
	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		return objectID.Hex()
	}

	return ""
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

func GetUserById(client *mongo.Client, id string) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("ID inválido:", err)
	}
	collection := client.Database("golang").Collection("users")
	filter := bson.D{{"_id", objectID}}
	var user bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User Found: %+v\n", user["name"])
}

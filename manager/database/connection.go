package database

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Collection {
	// Set Client Options
	clientOptions := options.Client().ApplyURI("mongodb+srv://ucmdls:ucmdls@ucmdls-manager.gk7um.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// get the collection
	collection := client.Database("ucmdls-manager").Collection("labs")

	return collection
}

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func Error(err error, w http.ResponseWriter) ErrorResponse {
	return ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
}

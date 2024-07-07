package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

type Event struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Date        string             `json:"date" bson:"date"`
	Location    string             `json:"location" bson:"location"`
	Capacity    int                `json:"capacity" bson:"capacity"`
}

type Response events.APIGatewayProxyResponse

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid ID",
		}, nil
	}

	var event Event

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://serverless-mongodb:27017"))
	collection := client.Database("eventdb").Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:      err.Error(),
		}, nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error encoding response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}

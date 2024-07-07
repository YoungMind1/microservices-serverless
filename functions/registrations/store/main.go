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
)

var client *mongo.Client

type Registration struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	EventID   string             `json:"event_id" bson:"event_id"`
	Canceled  bool               `json:"canceled" bson:"canceled"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}

type User struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

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
	var registration Registration
	if err := json.Unmarshal([]byte(req.Body), &registration); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	if registration.UserID == "" || registration.EventID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "user_id and event_id are required",
		}, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://serverless-mongodb:27017"))

	var user User
	collection := client.Database("userdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(registration.UserID)
    err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, nil
	}

	var event Event
	collection = client.Database("eventdb").Collection("events")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err = primitive.ObjectIDFromHex(registration.EventID)
    err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Event not found",
		}, nil
	}

 //    resp, err := http.Get("http://172.17.0.1:3000/api/users/" + url.QueryEscape(registration.UserID))
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 500,
	// 		Body:       err.Error(),
	// 	}, nil
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 400,
	// 		Body:       "This user_id does not exist",
	// 	}, nil
	// }
	//
 //    resp, err = http.Get("http://0.0.0.0:3000/api/events/" + url.QueryEscape(registration.EventID))
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 500,
	// 		Body:       err.Error(),
	// 	}, nil
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 400,
	// 		Body:       "This event_id does not exist",
	// 	}, nil
	// }

	collection = client.Database("registrationdb").Collection("registrations")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"user_id": registration.UserID, "event_id": registration.EventID, "canceled": false}).Err()
	if err != mongo.ErrNoDocuments {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "This user has already registered to this event",
		}, nil
	}

	registration.Canceled = false
	registration.Timestamp = time.Now()
	result, err := collection.InsertOne(ctx, registration)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	registration.ID = result.InsertedID.(primitive.ObjectID)
	body, err := json.Marshal(registration)
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

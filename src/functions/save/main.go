package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MySecret struct {
	MONGODB_URL string `json:"MONGODB_URL"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Create a new AWS session
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Secrets Manager client
	svc := secretsmanager.New(sess)

	fmt.Println(os.Getenv("SECRET_ARN"))

	// Retrieve the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("SECRET_ARN")),
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the secret value into a struct
	var secret MySecret
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB using the URL from the secret
	clientOptions := options.Client().ApplyURI(secret.MONGODB_URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get a handle to the "people" collection
	peopleCollection := client.Database("test-db").Collection("people")

	// Create a new person object
	person := Person{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	// Insert the person object into the "people" collection
	_, insertError := peopleCollection.InsertOne(ctx, person)
	if insertError != nil {
		log.Fatal(err)
	}

	// Return a success response
	return &events.APIGatewayProxyResponse{
		Body:       "Person saved to MongoDB",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

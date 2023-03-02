package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.mongodb.org/mongo-driver/bson"
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
	time1 := time.Now()
	fmt.Println(time1)

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

	time2 := time.Now()
	fmt.Println(time2)

	// Connect to MongoDB using the URL from the secret
	clientOptions := options.Client().ApplyURI(secret.MONGODB_URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	// Get a handle to the "people" collection
	peopleCollection := client.Database("test-db").Collection("people")

	time3 := time.Now()
	fmt.Println(time3)

	// Find all documents in the "people" collection
	cur, err := peopleCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	// defer cur.Close(ctx)

	time4 := time.Now()
	fmt.Println(time4)

	// Decode the documents into a slice of Person objects
	var people []Person
	err = cur.All(ctx, &people)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the slice of Person objects to a JSON array
	peopleBytes, err := json.Marshal(people)
	if err != nil {
		log.Fatal(err)
	}

	time5 := time.Now()
	fmt.Println(time5)

	// Return a success response
	return &events.APIGatewayProxyResponse{
		Body:       string(peopleBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

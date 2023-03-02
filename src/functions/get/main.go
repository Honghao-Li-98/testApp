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

var (
	client *mongo.Client
)

func init() {
	// Create a new AWS session
	sess, err1 := session.NewSession()
	if err1 != nil {
		log.Fatal(err1)
	}

	// Create a new Secrets Manager client
	svc := secretsmanager.New(sess)

	fmt.Println(os.Getenv("SECRET_ARN"))

	// Retrieve the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("SECRET_ARN")),
	}
	result, err2 := svc.GetSecretValue(input)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Parse the secret value into a struct
	var secret MySecret
	err3 := json.Unmarshal([]byte(*result.SecretString), &secret)
	if err3 != nil {
		log.Fatal(err3)
	}

	// Set up a MongoDB client
	clientOptions := options.Client().ApplyURI(secret.MONGODB_URL)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(timeoutContext context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Create a context with a timeout of 1 second
	timeoutContext, cancel := context.WithTimeout(timeoutContext, time.Second)
	defer cancel()

	// Get a handle to the "people" collection
	peopleCollection := client.Database("test-db").Collection("people")

	time3 := time.Now()
	fmt.Println(time3)

	// Find all documents in the "people" collection
	cur, err := peopleCollection.Find(timeoutContext, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	// defer cur.Close(timeoutContext)

	time4 := time.Now()
	fmt.Println(time4)

	// Decode the documents into a slice of Person objects
	var people []Person
	err = cur.All(timeoutContext, &people)
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

/*
2023-03-02T03:25:58.323-05:00	REPORT RequestId: 0d00d5d7-f29a-4993-826f-4ff207301ad4 Duration: 2.37 ms Billed Duration: 3 ms Memory Size: 128 MB Max Memory Used: 50 MB

2023-03-02T03:25:59.229-05:00	START RequestId: 33660246-3cef-4ed7-99cb-4587b4e5c6e9 Version: $LATEST

2023-03-02T03:25:59.229-05:00	2023-03-02 08:25:59.229783344 +0000 UTC m=+152.205445412

2023-03-02T03:25:59.231-05:00	2023-03-02 08:25:59.231066742 +0000 UTC m=+152.206728811

2023-03-02T03:25:59.231-05:00	2023-03-02 08:25:59.231098849 +0000 UTC m=+152.206760935

2023-03-02T03:25:59.231-05:00	END RequestId: 33660246-3cef-4ed7-99cb-4587b4e5c6e9

2023-03-02T03:25:59.231-05:00

Copy
REPORT RequestId: 33660246-3cef-4ed7-99cb-4587b4e5c6e9	Duration: 2.63 ms	Billed Duration: 3 ms	Memory Size: 128 MB	Max Memory Used: 50 MB
REPORT RequestId: 33660246-3cef-4ed7-99cb-4587b4e5c6e9 Duration: 2.63 ms Billed Duration: 3 ms Memory Size: 128 MB Max Memory Used: 50 MB
*/

func main() {
	lambda.Start(handler)
}

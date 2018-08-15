package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// question
type Question struct {
	ID       string    `json:"ID"`
	Category string `json:"Category"`
	Question string `json:"Question"`
	Answer1  string `json:"Answer1"`
	Answer2  string `json:"Answer2"`
	Answer3  string `json:"Answer3"`
	Answer4  string `json:"Answer4"`
}

const tableName = "Question"

// Create questions in database by http post
func PostQuestion(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}

	questionItem := &Question{}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	json.Unmarshal([]byte(request.Body), &questionItem)
	
	sensor, err := dynamodbattribute.MarshalMap(questionItem)

	// Create item in table question
	input := &dynamodb.PutItemInput{
		Item:      sensor,
		TableName: aws.String(tableName),
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}
	if _, err := svc.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:       "Bad Request," + err.Error(),
			StatusCode: 400,
		}, nil
	} else {
		return events.APIGatewayProxyResponse{ // success response
			Body:       "Created",
			StatusCode: 201,
		}, nil
	}
}

func main() {
	lambda.Start(PostQuestion)
}

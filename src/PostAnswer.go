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
	"github.com/satori/go.uuid"
)

// received data from wechat
type ReceiveData struct {
	Category1 string `json:"Category1"`
	Answer1   string `json:"Answer1"`
	Category2 string `json:"Category2"`
	Answer2   string `json:"Answer2"`
	Category3 string `json:"Category3"`
	Answer3   string `json:"Answer3"`
	Category4 string `json:"Category4"`
	Answer4   string `json:"Answer4"`
}

// Answer
type Answer struct {
	ID       string `json:"ID"`
	Category string `json:"Category"`
	Answer   string `json:"Answer"`
}

// Slice of Answer Data
type SliceAnsData struct {
	items []Answer
}

const tableName = "Answer"

// Create answers in database by http post
func PostAnswer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}

	receiveItem := &ReceiveData{}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	json.Unmarshal([]byte(request.Body), &receiveItem)

	items := []Answer{}
	queSlice := SliceAnsData{items}

	// get category1 and answer
	AnswerItem1 := Answer{}
	uid1, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	AnswerItem1.ID = uid1.String()
	AnswerItem1.Category = receiveItem.Category1
	AnswerItem1.Answer = receiveItem.Answer1
	queSlice.AddItem(AnswerItem1)
	// get category2 and answer
	AnswerItem2 := Answer{}
	uid2, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	AnswerItem2.ID = uid2.String()
	AnswerItem2.Category = receiveItem.Category2
	AnswerItem2.Answer = receiveItem.Answer2
	queSlice.AddItem(AnswerItem2)
	// get category3 and answer
	AnswerItem3 := Answer{}
	uid3, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	AnswerItem3.ID = uid3.String()
	AnswerItem3.Category = receiveItem.Category3
	AnswerItem3.Answer = receiveItem.Answer3
	queSlice.AddItem(AnswerItem3)
	// get category4 and answer
	AnswerItem4 := Answer{}
	uid4, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	AnswerItem4.ID = uid4.String()
	AnswerItem4.Category = receiveItem.Category4
	AnswerItem4.Answer = receiveItem.Answer4
	queSlice.AddItem(AnswerItem4)
	for _, i := range queSlice.items {
		sensor, err := dynamodbattribute.MarshalMap(i)

		// Create item in table Answer
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
	return events.APIGatewayProxyResponse{
		Body:       "Bad Request",
		StatusCode: 400,
	}, nil
}

// function to additem to slice
func (ansSlice *SliceAnsData) AddItem(item Answer) []Answer {
	ansSlice.items = append(ansSlice.items, item)
	return ansSlice.items
}

func main() {
	lambda.Start(PostAnswer)
}

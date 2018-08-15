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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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

// Slice of question Data
type SliceQueData struct {
	items []Question
}

// AWS region and DynamoDB table name
const region = "us-east-1"
const tableName = "Question"

var header = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
}

// to connect the database and to use get method to get the question with category from database
// https://url/question/{category}
// https://url/question/all for scan all the items in the table
func GetQuestion(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var category = request.PathParameters["category"]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    header,
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	var params *dynamodb.ScanInput
	if category == "all" {
		params = &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		}
	} else {
		flit := expression.Name("category").Equal(expression.Value(category))
		proj := expression.NamesList(expression.Name("ID"), expression.Name("Category"), expression.Name("Question"), expression.Name("Answer1"), expression.Name("Answer2"), expression.Name("Answer3"), expression.Name("Answer4"))
		expr, err := expression.NewBuilder().WithFilter(flit).WithProjection(proj).Build()

		if err != nil {
			panic(err)
		}
		params = &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(tableName),
		}
	}
	// scan all the item which meet requirments in the table
	result, err := svc.Scan(params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    header,
			Body:       "Not Found," + err.Error(),
			StatusCode: 404,
		}, nil
	}

	items := []Question{}
	queSlice := SliceQueData{items}
	// unmarshal the data from database
	for _, i := range result.Items {
		tempQuestion := Question{}
		err = dynamodbattribute.UnmarshalMap(i, &tempQuestion)
		queSlice.AddItem(tempQuestion)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Headers:    header,
				Body:       "Internal server error," + err.Error(),
				StatusCode: 500,
			}, nil
		}
	}
	outPut, err := json.Marshal(queSlice.items)
	if err != nil {
		panic(err)
	}
	return events.APIGatewayProxyResponse{
		Headers:    header,
		Body:       string(outPut),
		StatusCode: 200,
	}, nil
}

// function to additem to slice
func (queSlice *SliceQueData) AddItem(item Question) []Question {
	queSlice.items = append(queSlice.items, item)
	return queSlice.items
}

func main() {
	lambda.Start(GetQuestion)
}

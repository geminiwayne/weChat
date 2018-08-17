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

// AWS region and DynamoDB table name
const region = "us-east-1"
const tableName = "Answer"

var header = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
}

// to connect the database and to use get method to get the Answer with category from database
// https://url/answer/{category}
// https://url/answer/all for scan all the items in the table
func GetAnswer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
		proj := expression.NamesList(expression.Name("ID"), expression.Name("Category"),  expression.Name("Answer"))
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

	items := []Answer{}
	queSlice := SliceAnsData{items}
	// unmarshal the data from database
	for _, i := range result.Items {
		tempAnswer := Answer{}
		err = dynamodbattribute.UnmarshalMap(i, &tempAnswer)
		queSlice.AddItem(tempAnswer)
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
func (ansSlice *SliceAnsData) AddItem(item Answer) []Answer {
	ansSlice.items = append(ansSlice.items, item)
	return ansSlice.items
}

func main() {
	lambda.Start(GetAnswer)
}

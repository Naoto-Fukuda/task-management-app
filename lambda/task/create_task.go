package task

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	tableName = "TaskManagement"
	Svc       dynamodbiface.DynamoDBAPI
)

type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	DataType    string   `json:"DataType,omitempty"`
	DataValue   string   `json:"DataValue,omitempty"`
}

func AddTagToTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	taskId := request.QueryStringParameters["id"]
	tag := request.QueryStringParameters["tag"]

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(taskId),
			},
			"DataType": {
				S: aws.String("Tags"),
			},
		},
		UpdateExpression: aws.String("ADD dataValue :tag"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":tag": {
				SS: []*string{aws.String(tag)},
			},
		},
	}

	_, err := Svc.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to add tag to task: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Tag added to task successfully",
	}, nil
}

func CreateTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	task := Task{}
	err := json.Unmarshal([]byte(request.Body), &task)

	if err != nil {
		log.Fatalf("Failed to unmarshal task from JSON, %v", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// データタイプと値を設定
	if task.Title != "" {
		task.DataType = "Title"
		task.DataValue = task.Title
	} else if task.Status != "" {
		task.DataType = "Status"
		task.DataValue = task.Status
	} else if task.Description != "" {
		task.DataType = "Description"
		task.DataValue = task.Description
	} else if len(task.Tags) > 0 {
		task.DataType = "Tags"
		tags, err := json.Marshal(task.Tags)
		if err != nil {
			log.Fatalf("Failed to marshal tags into JSON, %v", err)
			return events.APIGatewayProxyResponse{}, err
		}
		task.DataValue = string(tags)
	}

	av, err := dynamodbattribute.MarshalMap(task) // タスクオブジェクトをDynamoDBのアイテムに変換
	if err != nil {
		log.Fatalf("Failed to marshal task into DynamoDB item, %v", err)
		return events.APIGatewayProxyResponse{}, err
	}
	if _, ok := av["DataType"]; !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing DataType in the item",
		}, fmt.Errorf("missing DataType in the item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = Svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "Task created successfully",
		StatusCode: http.StatusCreated,
	}, nil
}

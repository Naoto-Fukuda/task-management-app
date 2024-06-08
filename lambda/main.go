package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	tableName = "TaskManagement"
	// svc  = mockDynamoDBClient{}
	// svc       *dynamodb.DynamoDB
	svc       dynamodbiface.DynamoDBAPI
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

func getTaskById(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.QueryStringParameters["id"]

	// Query APIを使用して複数のアイテムを取得
	input := &dynamodb.QueryInput{
			TableName: aws.String(tableName),
			KeyConditionExpression: aws.String("id = :id"), // Where句のパラメータ
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ 
					":id": {
							S: aws.String(id),
					},
			},
	}

	// DynamoDBクライアントの呼び出し
	result, err := svc.Query(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	
	taskMap := make(map[string]*Task)

	for _, i := range result.Items {
		id := *i["id"].S
		if _, exists := taskMap[id]; !exists {
			taskMap[id] = &Task{ID: id}
		}
		task := taskMap[id]

		switch *i["dataType"].S {
		case "Title":
			task.Title = *i["dataValue"].S
		case "Description":
			task.Description = *i["dataValue"].S
		case "Status":
			task.Status = *i["dataValue"].S
		case "Tags":
			if task.Tags == nil {
				task.Tags = []string{}
			}
			task.Tags = append(task.Tags, *i["dataValue"].S)
		}
	}

	// 結果をリストに変換
	tasks := make([]Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

// func getTasks(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

// }
//
//	func getTasksByTitle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		title := request.QueryStringParameters["title"]
//	}
//
//	func getTasksByDescription(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		description := request.QueryStringParameters["description"]
//	}
//
//	func getTasksByStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		status := request.QueryStringParameters["status"]
//	}
//
//	func getTasksByTag(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		tag := request.QueryStringParameters["tag"]
//	}
func createTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "Task created successfully",
		StatusCode: http.StatusCreated,
	}, nil
}

//	func addTagToTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		taskId := request.QueryStringParameters["id"]
//		tag := request.QueryStringParameters["tag"]
//	}
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := request.HTTPMethod
	switch method {
	case "GET":
		return getTaskById(request)
	case "POST":
		return createTask(request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "Method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

func main() {
	svc = dynamodb.New(session.Must(session.NewSession()))
	lambda.Start(handler)
}

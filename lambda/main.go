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

func getTaskById(id string) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
		},
	}

	result, err := svc.Query(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Query failed: %v", err),
		}, nil
	}

	// タスクの地図を作成
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

	tasks := make([]Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to marshal tasks: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

// TODO:ButchGetについて整理する
// あと、responseがmapなのでそれの扱い方を再度確認する
//　keyで取得するとか更新するとか
func getTasksByTaskIds(ids []string) (map[string]*Task, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(ids))
	for i, id := range ids {
		keys[i] = map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		}
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: keys,
			},
		},
	}

	result, err := svc.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	taskMap := make(map[string]*Task)
	for _, i := range result.Responses[tableName] {
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

	return taskMap, nil
}

func getTasksByTitle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	title := request.QueryStringParameters["title"]

	// Titleに基づいてタスクIDを取得
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("dataType = :dataType AND dataValue = :dataValue"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dataType": {
				S: aws.String("Title"),
			},
			":dataValue": {
				S: aws.String(title),
			},
		},
	}

	result, err := svc.Query(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Query failed: %v", err),
		}, nil
	}

	idsMap := make(map[string]bool)
	for _, i := range result.Items {
		id := *i["id"].S
		if _, exists := idsMap[id]; !exists {
			idsMap[id] = true
		}
	}

	ids := make([]string, 0, len(idsMap))
	for id := range idsMap {
		ids = append(ids, id)
	}

	taskMap, err := getTasksByTaskIds(ids)
	if err != nil {
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       fmt.Sprintf("Failed to retrieve tasks by IDs: %v", err),
    }, nil
	}

	tasks := make([]Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to marshal response: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

	func getTasksByStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		status := request.QueryStringParameters["status"]

		input := &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String("GSI1"),
			KeyConditionExpression: aws.String("dataType = :dataType AND dataValue = :dataValue"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":dataType": {
					S: aws.String("Status"),
				},
				":dataValue": {
					S: aws.String(status),
				},
			},
		}

			result, err := svc.Query(input)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       fmt.Sprintf("Query failed: %v", err),
				}, nil
			}

			idsMap := make(map[string]bool)
			for _, i := range result.Items {
				id := *i["id"].S
				if _, exists := idsMap[id]; !exists {
					idsMap[id] = true
				}
			}

			ids := make([]string, 0, len(idsMap))
			for id := range idsMap {
				ids = append(ids, id)
			}

			taskMap, err := getTasksByTaskIds(ids)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       fmt.Sprintf("Failed to retrieve tasks by IDs: %v", err),
				}, nil
			}

			tasks := make([]Task, 0, len(taskMap))
			for _, task := range taskMap {
				tasks = append(tasks, *task)
			}
			
			response, err := json.Marshal(tasks)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       fmt.Sprintf("Failed to marshal response: %v", err),
				}, nil
			}

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(response),
			}, nil
	}

	func getTasksByTag(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tag := request.QueryStringParameters["tag"]

		input := &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String("GSI1"),
			KeyConditionExpression: aws.String("dataType = :dataType AND dataValue = :dataValue"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":dataType": {
					S: aws.String("Tags"),
				},
				":dataValue": {
					S: aws.String(tag),
				},
			},
		}

		result, err := svc.Query(input)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Query failed: %v", err),
			}, nil
		}

		idsMap := make(map[string]bool)
		for _, i := range result.Items {
			id := *i["id"].S
			if _, exists := idsMap[id]; !exists {
				idsMap[id] = true
			}
		}

		ids := make([]string, 0, len(idsMap))
		for id := range idsMap {
			ids = append(ids, id)
		}

		taskMap, err := getTasksByTaskIds(ids)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Failed to retrieve tasks by IDs: %v", err),
			}, nil
		}

		tasks := make([]Task, 0, len(taskMap))
		for _, task := range taskMap {
			tasks = append(tasks, *task)
		}

		response, err := json.Marshal(tasks)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Failed to marshal response: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(response),
		}, nil
}

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

// これはupdateやdeleteと共通化できるかも？
func addTagToTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	_, err := svc.UpdateItem(input)
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

// updateは共通関数で引数を変えるだけでいけるかも?
// func updateTagOnTask
// updateTitleOnTask
// updateDescriptionOnTask
// func deleteTaskById

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := request.HTTPMethod
	switch method {
	case "GET":
		return getTaskById(request.PathParameters["id"])
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

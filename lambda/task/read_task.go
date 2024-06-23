package task

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetTaskById(id string) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
		},
	}

	result, err := Svc.Query(input)
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
		dataType := *i["dataType"].S
		dataValue := *i["dataValue"].S
		UpdateTaskField(task, dataType, dataValue)
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

func GetTasksByTaskIds(ids []string) (map[string]*Task, error) {
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

	result, err := Svc.BatchGetItem(input)
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

func GetTasksByAttribute(request events.APIGatewayProxyRequest, attributeKey string, attributeValue string) (events.APIGatewayProxyResponse, error) {
	value := request.QueryStringParameters[attributeKey]

	input := &dynamodb.QueryInput{
		TableName:              aws.String("TaskManagement"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("dataType = :dataType AND dataValue = :dataValue"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dataType": {
				S: aws.String(attributeKey),
			},
			":dataValue": {
				S: aws.String(value),
			},
		},
	}

	result, err := Svc.Query(input)
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

	taskMap, err := GetTasksByTaskIds(ids)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to retrieve tasks by %v", attributeKey),
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

func GetTasksByTag(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	result, err := Svc.Query(input)
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

	taskMap, err := GetTasksByTaskIds(ids)
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
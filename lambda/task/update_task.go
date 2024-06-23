package task

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func UpdateTaskField(task *Task, dataType string, dataValue string) {
	switch dataType {
	case "Title":
		task.Title = dataValue
	case "Description":
		task.Description = dataValue
	case "Status":
		task.Status = dataValue
	case "Tags":
		if task.Tags == nil {
			task.Tags = []string{}
		}
		task.Tags = append(task.Tags, dataValue)
	}
}

func UpdateTagOnTask(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	taskId := request.QueryStringParameters["id"]
	old_tag := request.QueryStringParameters["old_tag"]
	new_tag := request.QueryStringParameters["new_tag"]

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
		UpdateExpression: aws.String("DELETE dataValue :old_tag ADD dataValue :new_tag"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":old_tag": {
				SS: []*string{aws.String(old_tag)},
			},
			":new_tag": {
				SS: []*string{aws.String(new_tag)},
			},
		},
	}

	_, err := Svc.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to update tag on task: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Tag updated on task successfully",
	}, nil
}

func UpdateTaskAttribute(request events.APIGatewayProxyRequest, attributeKey string, attributeValue string) (events.APIGatewayProxyResponse, error) {
	taskId := request.QueryStringParameters["id"]
	newValue := request.QueryStringParameters[attributeValue]

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(taskId),
			},
			"DataType": {
				S: aws.String(attributeKey),
			},
		},
		UpdateExpression: aws.String("SET dataValue = :new_value"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":new_value": {
				S: aws.String(newValue),
			},
		},
	}

	_, err := Svc.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to update %s on task: %v", attributeKey, err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%s updated on task successfully", attributeKey),
	}, nil
}
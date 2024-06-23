package task

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DeleteTaskById(id string) (response events.APIGatewayProxyResponse, err error) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err = Svc.DeleteItem(input)
	if err != nil {
		response = events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to delete task: %v", err),
		}
		return
	}

	response = events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Task deleted successfully",
	}
	return
}

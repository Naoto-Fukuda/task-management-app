package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/mock"
)

// DynamoDBAPIのモック定義
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

// PutItemのモック実装
func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

// GetItemのモック実装
func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

// DeleteItemのモック実装
func (m *mockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

// Queryのモック実装
func (m *mockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

// BatchGetItemのモック実装
func (m *mockDynamoDBClient) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.BatchGetItemOutput), args.Error(1)
}

// UpdateItemのモック実装
func (m *mockDynamoDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

// func setup() {
// 	// 本番環境では実際のDynamoDBクライアントを使用する
// 	svc = &mockDynamoDBClient{}
// }

// var svc *dynamodb.DynamoDB

// func teardown() {
// 	// クリーンアップ処理が必要な場合はここに記載
// 	svc = nil
// }

func Test_createTask(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	svc = mockSvc

	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Request",
			args: args{
				request: events.APIGatewayProxyRequest{
					Body:       "{\"id\":\"1\", \"dataType\":\"Description\", \"dataValue\":\"Task Title\"}",
					HTTPMethod: "POST",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "Task created successfully",
				StatusCode: http.StatusCreated,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTask(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTaskById(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":        {S: aws.String("1")},
				"dataType":  {S: aws.String("Title")},
				"dataValue": {S: aws.String("Task Title")},
			},
			{
				"id":        {S: aws.String("1")},
				"dataType":  {S: aws.String("Description")},
				"dataValue": {S: aws.String("Description of the task1")},
			},
			{
				"id":        {S: aws.String("2")},
				"dataType":  {S: aws.String("Title")},
				"dataValue": {S: aws.String("Task Title")},
			},
			{
				"id":        {S: aws.String("2")},
				"dataType":  {S: aws.String("Description")},
				"dataValue": {S: aws.String("Description of the task2")},
			},
		},
	}, nil)

	svc = mockSvc

	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid ID",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{"id": "1"},
					HTTPMethod:            "GET",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "[{\"id\":\"1\",\"title\":\"Task Title\",\"description\":\"Description of the task1\"},{\"id\":\"2\",\"title\":\"Task Title\",\"description\":\"Description of the task2\"}]",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskID := tt.args.request.QueryStringParameters["id"]
			got, err := getTaskById(taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTasksById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTasksById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTasksByTitle(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":        {S: aws.String("1")},
				"dataType":  {S: aws.String("Title")},
				"dataValue": {S: aws.String("Task Title")},
			},
			{
				"id":        {S: aws.String("2")},
				"dataType":  {S: aws.String("Title")},
				"dataValue": {S: aws.String("Task Title")},
			},
		},
	}, nil)
	mockSvc.On("BatchGetItem", mock.Anything).Return(&dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			"TaskManagement": {
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag1")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task1")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag2")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task2")},
				},
			},
		},
	}, nil)

	svc = mockSvc
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Title",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{"title": "Task Title"},
					HTTPMethod:            "GET",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "[{\"id\":\"1\",\"title\":\"Task Title\",\"description\":\"Description of the task1\",\"tags\":[\"Tag1\"]},{\"id\":\"2\",\"title\":\"Task Title\",\"description\":\"Description of the task2\",\"tags\":[\"Tag2\"]}]",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTasksByTitle(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTasksByTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTasksByTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTasksByStatus(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":        {S: aws.String("1")},
				"dataType":  {S: aws.String("Status")},
				"dataValue": {S: aws.String("Completed")},
			},
			{
				"id":        {S: aws.String("2")},
				"dataType":  {S: aws.String("Status")},
				"dataValue": {S: aws.String("Completed")},
			},
		},
	}, nil)
	mockSvc.On("BatchGetItem", mock.Anything).Return(&dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			"TaskManagement": {
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag1")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task1")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag2")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task2")},
				},
			},
		},
	}, nil)

	svc = mockSvc

	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Status",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{"status": "Completed"},
					HTTPMethod:            "GET",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "[{\"id\":\"1\",\"title\":\"Task Title\",\"description\":\"Description of the task1\",\"tags\":[\"Tag1\"]},{\"id\":\"2\",\"title\":\"Task Title\",\"description\":\"Description of the task2\",\"tags\":[\"Tag2\"]}]",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTasksByStatus(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTasksByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTasksByStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTasksByTag(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":        {S: aws.String("1")},
				"dataType":  {S: aws.String("Tags")},
				"dataValue": {S: aws.String("Tag1")},
			},
			{
				"id":        {S: aws.String("2")},
				"dataType":  {S: aws.String("Tags")},
				"dataValue": {S: aws.String("Tag1")},
			},
		},
	}, nil)

	mockSvc.On("BatchGetItem", mock.Anything).Return(&dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			"TaskManagement": {
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag1")},
				},
				{
					"id":        {S: aws.String("1")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task1")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Title")},
					"dataValue": {S: aws.String("Task Title")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Tags")},
					"dataValue": {S: aws.String("Tag1")},
				},
				{
					"id":        {S: aws.String("2")},
					"dataType":  {S: aws.String("Description")},
					"dataValue": {S: aws.String("Description of the task2")},
				},
			},
		},
	}, nil)

	svc = mockSvc
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Tag",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{"tag": "Tag1"},
					HTTPMethod:            "GET",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "[{\"id\":\"1\",\"title\":\"Task Title\",\"description\":\"Description of the task1\",\"tags\":[\"Tag1\"]},{\"id\":\"2\",\"title\":\"Task Title\",\"description\":\"Description of the task2\",\"tags\":[\"Tag1\"]}]",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTasksByTag(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTasksByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTasksByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addTagToTask(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("UpdateItem", mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	svc = mockSvc
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Request",
			args: args{
				request: events.APIGatewayProxyRequest{
					Body:       "{\"id\":\"1\", \"tag\":\"Tag1\"}",
					HTTPMethod: "POST",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "Tag added to task successfully",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addTagToTask(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("addTagToTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addTagToTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteTaskById(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

	svc = mockSvc
	type args struct {
		id string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Valid ID",
			args: args{
				id: "1",
			},
			wantResponse: events.APIGatewayProxyResponse{
				Body:       "Task deleted successfully",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := deleteTaskById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteTaskById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("deleteTaskById() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_updateTagOnTask(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("UpdateItem", mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	svc = mockSvc
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Request",
			args: args{
				request: events.APIGatewayProxyRequest{
					Body:       "{\"id\":\"1\", \"oldTag\":\"Tag1\", \"newTag\":\"Tag2\"}",
					HTTPMethod: "PUT",
				},
			},
			want: events.APIGatewayProxyResponse{
				Body:       "Tag updated on task successfully",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updateTagOnTask(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateTagOnTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateTagOnTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateTaskAttribute(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	mockSvc.On("UpdateItem", mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	svc = mockSvc
	type args struct {
		request        events.APIGatewayProxyRequest
		attributeKey   string
		attributeValue string
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Valid Request",
			args: args{
				request: events.APIGatewayProxyRequest{
					Body:       "{\"id\":\"1\", \"attributeValue\":\"Completed\"}",
					HTTPMethod: "PUT",
				},
				attributeKey:   "Status",
				attributeValue: "Completed",
			},
			want: events.APIGatewayProxyResponse{
				Body:       "Status updated on task successfully",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updateTaskAttribute(tt.args.request, tt.args.attributeKey, tt.args.attributeValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateTaskAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateTaskAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

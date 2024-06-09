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

// Queryのモック実装
func (m *mockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
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

func Test_getTasksById(t *testing.T) {
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
				"dataValue": {S: aws.String("Description of the task")},
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
				Body:       "[{\"id\":\"1\",\"title\":\"Task Title\",\"description\":\"Description of the task\"}]",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTaskById(tt.args.request)
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
				},
				{
					"id":        {S: aws.String("2")},
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

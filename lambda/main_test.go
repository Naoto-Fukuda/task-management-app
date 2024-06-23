package main

import (
	"net/http"
	"reflect"
	"testing"
	"task-management-app/lambda/task"
	"task-management-app/lambda/mocks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/golang/mock/gomock"
)

// DynamoDBAPIのモック定義
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	ctrl *gomock.Controller
}

// モックオブジェクトの生成
func newMockDynamoDBClient(ctrl *gomock.Controller) *mockDynamoDBClient {
	mock := &mockDynamoDBClient{ctrl: ctrl}
	return mock
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// DynamoDBのモッククライアントを作成
	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	mockDynamoDB.EXPECT().PutItem(gomock.Any()).Return(&dynamodb.PutItemOutput{}, nil).Times(1)

	task.Svc = mockDynamoDB

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
			got, err := task.CreateTask(tt.args.request)
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

func Test_GetTaskById(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// DynamoDBのモッククライアントを作成
	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	mockDynamoDB.EXPECT().Query(gomock.Any()).Return(&dynamodb.QueryOutput{
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
	}, nil).Times(1) // 期待される呼び出し回数を指定

	task.Svc = mockDynamoDB

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
			got, err := task.GetTaskById(taskID)
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

func Test_getTasksByTag(t *testing.T) {
	// DynamoDBのモッククライアントを作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	mockDynamoDB.EXPECT().Query(gomock.Any()).Return(&dynamodb.QueryOutput{
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
	}, nil).Times(1)

	mockDynamoDB.EXPECT().BatchGetItem(gomock.Any()).Return(&dynamodb.BatchGetItemOutput{
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
	}, nil).Times(1)

	task.Svc = mockDynamoDB

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
			got, err := task.GetTasksByTag(tt.args.request)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	mockDynamoDB.EXPECT().UpdateItem(gomock.Any()).Return(&dynamodb.UpdateItemOutput{}, nil).Times(1)

	task.Svc = mockDynamoDB
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
			got, err := task.AddTagToTask(tt.args.request)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	task.Svc = mockDynamoDB

	// DeleteItemの期待される呼び出しを設定
	mockDynamoDB.EXPECT().DeleteItem(gomock.Any()).Return(&dynamodb.DeleteItemOutput{}, nil).Times(1)

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
			gotResponse, err := task.DeleteTaskById(tt.args.id)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	task.Svc = mockDynamoDB
	mockDynamoDB.EXPECT().UpdateItem(gomock.Any()).Return(&dynamodb.UpdateItemOutput{}, nil).Times(1)

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
			got, err := task.UpdateTagOnTask(tt.args.request)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	mockDynamoDB.EXPECT().UpdateItem(gomock.Any()).Return(&dynamodb.UpdateItemOutput{}, nil).Times(1)

	task.Svc = mockDynamoDB
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
			got, err := task.UpdateTaskAttribute(tt.args.request, tt.args.attributeKey, tt.args.attributeValue)
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

func Test_getTasksByAttribute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDynamoDB := mockdb.NewMockDynamoDBAPI(ctrl)
	task.Svc = mockDynamoDB

	// Queryの期待される呼び出しを設定（引数を具体的に指定）
	mockDynamoDB.EXPECT().Query(&dynamodb.QueryInput{
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
        ":dataValue": {S: aws.String("Task Title")},
        ":dataType": {S: aws.String("Title")},
    },
    IndexName: aws.String("GSI1"),
    KeyConditionExpression: aws.String("dataType = :dataType AND dataValue = :dataValue"),
    TableName: aws.String("TaskManagement"),
}).Return(&dynamodb.QueryOutput{
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


	mockDynamoDB.EXPECT().BatchGetItem(gomock.Any()).Return(&dynamodb.BatchGetItemOutput{
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
			name: "Valid Title",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{"Title": "Task Title"},
					HTTPMethod:            "GET",
				},
				attributeKey:   "Title",
				attributeValue: "Task Title",
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
			got, err := task.GetTasksByAttribute(tt.args.request, tt.args.attributeKey, tt.args.attributeValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)


func NewStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	readCapacity := float64(2)
	writeCapacity := float64(2)

	stack := awscdk.NewStack(scope, &id, props)
	table := awsdynamodb.NewTable(stack, jsii.String("TaskManagement"), &awsdynamodb.TableProps{
		TableName: jsii.String("TaskManagement"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("DataType"),
			Type: awsdynamodb.AttributeType_STRING},
		DeletionProtection: aws.Bool(true),
		ReadCapacity:       &readCapacity,
		WriteCapacity:      &writeCapacity,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("GSI1"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("DataValue"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ReadCapacity:  &readCapacity,
		WriteCapacity: &writeCapacity,
	})

	// Create a Lambda function
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("TaskManagementFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Code:    awslambda.Code_FromAsset(jsii.String("../lambda"), nil),
		Handler: jsii.String("bootstrap"),
	})
	
		// Grant the Lambda function read/write permissions to the table
	table.GrantReadWriteData(lambdaFunction)

	return stack
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := awscdk.NewApp(nil)

	NewStack(app, "taskManagementStack", &awscdk.StackProps{
			Env: env(),
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	awsAccountId := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_REGION")

	return &awscdk.Environment{
		Account: jsii.String(awsAccountId),
		Region:  jsii.String(awsRegion),
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dorneanu/aws-custom-resource-poc/internal"
)

var awsSession aws.Config

// init will setup the AWS session
func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	awsSession = cfg
	fmt.Print("initialized aws session: %#v", awsSession)
}

// lambdaHandler handles incoming CloudFormation events
// is of type cfn.CustomResourceFunction
func lambdaHandler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	switch event.ResourceType {
		cfn.Custom
	case "Custom::SSMCredential":
		resourceHandler := internal.NewSSMCustomResourceHandler(awsSession)

	default:
		return "", fmt.Errorf("Unknown resource type: %s", event.ResourceType)
	}
}

// main function
func main() {
	// From : https://github.com/aws/aws-lambda-go/blob/main/cfn/wrap.go
	//
	// LambdaWrap returns a CustomResourceLambdaFunction which is something lambda.Start()
	// will understand. The purpose of doing this is so that Response Handling boiler
	// plate is taken away from the customer and it makes writing a Custom Resource
	// simpler.
	//
	//	func myLambda(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	//		physicalResourceID = "arn:...."
	//		return
	//	}
	//
	//	func main() {
	//		lambda.Start(cfn.LambdaWrap(myLambda))
	//	}
	lambda.Start(cfn.LambdaWrap(lambdaHandler))
}

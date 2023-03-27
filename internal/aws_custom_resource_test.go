package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMParameterApiImpl is a mock for SSMParameterAPI
type SSMParameterApiImpl struct{}

// PutParameter ...
// TODO: Implement this
func (s SSMParameterApiImpl) PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	output := &ssm.PutParameterOutput{}
	return output, nil
}

// DeleteParameter ...
func (s SSMParameterApiImpl) DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error) {
	output := &ssm.DeleteParameterOutput{}
	return output, nil
}

// TestPutParameter ...
func TestPutParameter(t *testing.T) {
	fmt.Printf("Implement this ")
	mockedAPI := SSMParameterApiImpl{}
	ssmHandler := SSMCustomResourceHandler{
		ssmClient: mockedAPI,
	}

	// Create new SSM parameter
	cfnEvent := cfn.Event{
		RequestType:        "Create",
		RequestID:          "xxx",
		ResponseURL:        "https://dornea.nu",
		ResourceType:       "AWS::CloudFormation::CustomResource",
		PhysicalResourceID: "",
		LogicalResourceID:  "SSMCredentialTesting1",
		StackID:            "arn:aws:cloudformation:eu-central-1:9999999:stack/CustomResourcesGolang",
		ResourceProperties: map[string]interface{}{
			"ServiceToken": "arn:aws:lambda:eu-central-1:9999999:function:CustomResourcesGolang-Function",
			"key":          "/testing3",
			"value":        "some-secret-value",
		},
	}
	_, _, _ = ssmHandler.Create(context.TODO(), cfnEvent)

}

// TestDeleteParameter ...
// TODO: Implement this
func TestDeleteParameter(t *testing.T) {
	fmt.Printf("Implement this")
}

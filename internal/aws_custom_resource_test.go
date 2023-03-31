package internal

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMParameterApiImpl is a mock for SSMParameterAPI
type SSMParameterApiImpl struct{}

// PutParameter
func (s SSMParameterApiImpl) PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	output := &ssm.PutParameterOutput{}
	return output, nil
}

// DeleteParameter
func (s SSMParameterApiImpl) DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error) {
	output := &ssm.DeleteParameterOutput{}
	return output, nil
}

// TestDeleteParameter
func TestDeleteParameter(t *testing.T) {
	// mockAPI := SSMParameterApiImpl{}
	// ssmHandler := SSMCustomResourceHandler{
	// 	ssmClient: mockAPI,
	// }
	// cfnEvent := cfn.Event{
	// 	RequestType:        "Delete",
	// 	RequestID:          "xxx",
	// 	ResponseURL:        "https://dornea.nu",
	// 	ResourceType:       "AWS::CloudFormation::CustomResource",
	// 	PhysicalResourceID: "arn:aws:ssm:eu-central-1:9999999:parameter/testing3",
	// 	LogicalResourceID:  "SSMCredentialTesting1",
	// }

}

// TestPutParameter ...
func TestPutParameter(t *testing.T) {
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

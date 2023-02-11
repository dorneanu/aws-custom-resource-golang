package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSMParameterAPI defines an interface for the SSM API calls
// I use this interface in order to be able to mock out the SSM
// client and implement unit tests properly.
// Also check https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2/ssm
type SSMParameterAPI interface {
	DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

type SSMCustomResourceHandler struct {
	ssmClient SSMParameterAPI
}

// NewSSMCustomResourceHandler returns a new handler for the SSM custom resource
func NewSSMCustomResourceHandler(cfg aws.Config) SSMCustomResourceHandler {
	return SSMCustomResourceHandler{
		ssmClient: ssm.NewFromConfig(cfg),
	}
}

// handle does something
func (s SSMCustomResourceHandler) handleSSMCustomResource(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
	switch event.RequestType {
	case "Create":
		return s.Create(ctx, event)
	case "Update":
		return s.Update(ctx, event)
	case "Delete":
		return s.Delete(ctx, event)
	default:
		return "", fmt.Errorf("Unknown request type: %s", event.RequestType)
	}
}

// Create puts a new SSM parameter
func (s SSMCustomResourceHandler) Create(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
	var physicalResourceID string

	// Get custom resource parameter from event
	ssmPath, err := strProperty(event, "credential_name")
	if err != nil {
		return physicalResourceID, fmt.Errorf("Couldn't extract credential_name: %s", err)
	}
	physicalResourceID = ssmPath

	ssmValue, err := strProperty(event, "credential_value")
	if err != nil {
		return physicalResourceID, fmt.Errorf("Couldn't extract credential_value: %s", err)
	}

	// Put new parameter
	_, err = s.ssmClient.PutParameter(context.Background(), &ssm.PutParameterInput{
		Name:      aws.String(ssmPath),
		Value:     aws.String(ssmValue),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(true),
	})

	if err != nil {
		return physicalResourceID, fmt.Errorf("Couldn't put parameter (%s): %s\n", ssmPath, err)
	}
	return physicalResourceID, nil
}

// Update overwrites a SSM parameter by a new value
func (s SSMCustomResourceHandler) Update(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
	return s.Create(ctx, event)
}

// Delete will delete a SSM parameter
func (s SSMCustomResourceHandler) Delete(ctx context.Context, event cfn.Event) (string, map[string]interface{}, error) {
	var physicalResourceID string

	ssmPath, err := strProperty(event, "credential_name")
	if err != nil {
		return physicalResourceID, fmt.Errorf("Couldn't find property credential_name: %s", err)
	}
	physicalResourceID = ssmPath

	_, err = s.ssmClient.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{
		Name: aws.String(ssmPath),
	})

	if err != nil {
		return physicalResourceID, fmt.Errorf("Couldn't delete parameter (%s): %s\n", ssmPath, err)
	}

	return physicalResourceID, nil
}

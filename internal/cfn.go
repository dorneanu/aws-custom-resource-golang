package internal

import (
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
)

// strProperty extracts a property name from a CFN event
// From https://github.com/masgari/aws-custom-resources/blob/d2705f3c8c8bc2fa39b43b5da60fb998a8c660e7/internal/common.go#L9
func strProperty(event cfn.Event, propertyName string) (string, error) {
	if val, ok := event.ResourceProperties[propertyName]; ok {
		return val.(string), nil
	}
	return "", fmt.Errorf("Missing property %s", propertyName)
}

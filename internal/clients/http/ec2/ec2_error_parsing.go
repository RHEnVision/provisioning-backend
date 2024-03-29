package ec2

import (
	"errors"
	"strings"

	"github.com/aws/smithy-go"
)

func isAWSUnauthorizedError(err error) bool {
	return isAWSOperationError(err, "api error UnauthorizedOperation")
}

func isAWSOperationError(err error, substr string) bool {
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		return strings.Contains(oe.Unwrap().Error(), substr)
	}
	return false
}

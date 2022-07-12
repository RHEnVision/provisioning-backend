package ec2

import (
	"errors"
	"strings"

	"github.com/aws/smithy-go"
)

var DuplicatePubkeyErr = errors.New("pubkey already exists")
var OperationNotPermittedErr = errors.New("operation not permitted")
var MoreThan100InstanceTypes = errors.New("there are more then 100 instance types")

func IsOperationError(err error, substr string) bool {
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) && strings.Contains(oe.Unwrap().Error(), substr) {
			return true
		}
	}
	return false
}

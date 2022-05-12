package ec2

import (
	"errors"
	"strings"

	"github.com/aws/smithy-go"
)

var DuplicatePubkeyErr = errors.New("pubkey already exists")
var OperationNotPermittedErr = errors.New("operation not permitted")

func IsOperationError(err error, substr string) bool {
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) && strings.Contains(oe.Unwrap().Error(), substr) {
			return true
		}
	}
	return false
}

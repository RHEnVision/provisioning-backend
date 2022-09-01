package stubs

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/smithy-go/ptr"
)

type stsCtxKeyType string

var stsCtxKey stsCtxKeyType = "sts-interface"

type STSClientStub struct{}

func init() {
	clients.GetSTSClient = getSTSClientStub
}

// STSClient
func WithSTSClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, stsCtxKey, &STSClientStub{})
	return ctx
}

func getSTSClientStub(ctx context.Context) (si clients.STS, err error) {
	var ok bool
	if si, ok = ctx.Value(stsCtxKey).(*STSClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *STSClientStub) AssumeRole(arn string) (*stsTypes.Credentials, error) {
	d := time.Date(2023, 11, 17, 20, 34, 58, 651387237, time.UTC)
	return &stsTypes.Credentials{
		AccessKeyId:     ptr.String("ASIATLQ4WO467P5JAAAD"),
		SecretAccessKey: ptr.String("3Y7LGRfxoKwZwwFg9Rvp9rtRxWnK2UwC5kTSc0vy"),
		Expiration:      &d,
		SessionToken:    ptr.String("IQoJb3JpZ2luX2VjENX//////////wEaCXVzLWVhc3QtMSJHMEUCIEdy/horMGYAwgz+gs7sHuTuMdWq77kZxpKQ3D/Bn3kDAiEAsIHxD5V9UhuV+k1Kgh2wfgXeK/tQgn1DCf4+BmMAL9cqkQIIXRAAGgwyMzA5MTQ2ODQ3MzMiDD4tHpfCWvKHO/je5iruAfY8hoJZRQ2N+wEQsl5bWBaND9UvZaJ9cSza8yCkw7ACfmB5KvZ6fHEUeNMNXTRkxE4SVtdXvyGFXPn8WmV2SthvUw/SwaVzIG/z4qTU58qV/G14TKXqFmFYN6lAOvPmzSnWHm6b5m2b+CMHujG7BDUkeQ9zvmqgf597L4NFd/43a5DLInEyhHRiTKYO89ifWFcik8uGw8vssalPipa3eD5NvifuMBig45Qj0n+tkrTd7qPX1D+QPQvNoCtIoexr1pP7NEIUOVJSBoh+oEeOXLq0OgL6CEZ3i6ouHw/TGI/KQvsZhyfIxmasTGxTmKEwjqS9mAY6nQGhzef34qVsZN1ZThWIYf3BYco+gt5DCscYSp0Uxo4ofVMFZlAKI4O/hVyUksBhxW4sjyYVgRsyMoX7//brYFc3ulioorI5RsxkYZaYifzdWAvwW1tuJIqJhSfZM/o0Bhd2W5BmR2ZocJye7iDSk6Eqdx7h4xv55tgIIIL6BHH+5rIK/Zn1Q5Hw/qDdN3R8+/l90xUErdQXFLV/zP8v"),
	}, nil
}

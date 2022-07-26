package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type imageBuilderCtxKeyType string

var imageBuilderCtxKey imageBuilderCtxKeyType = "image-builder-interface"

type ImageBuilderClientStub struct{}

func init() {
	clients.GetImageBuilderClient = getImageBuilderClientStub
}

type contextReadError struct{}

func (m *contextReadError) Error() string {
	return "failed to find or convert dao stored in testing context"
}

func WithImageBuilderClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, imageBuilderCtxKey, &ImageBuilderClientStub{})
	return ctx
}

func getImageBuilderClientStub(ctx context.Context) (si clients.ImageBuilder, err error) {
	var ok bool
	if si, ok = ctx.Value(imageBuilderCtxKey).(*ImageBuilderClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}
func (mock *ImageBuilderClientStub) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	return "ami-0c830793775595d4b-test", nil
}

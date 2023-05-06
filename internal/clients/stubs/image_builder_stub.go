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

func WithImageBuilderClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, imageBuilderCtxKey, &ImageBuilderClientStub{})
	return ctx
}

func getImageBuilderClientStub(ctx context.Context) (si clients.ImageBuilder, err error) {
	var ok bool
	if si, ok = ctx.Value(imageBuilderCtxKey).(*ImageBuilderClientStub); !ok {
		err = ContextReadError
	}
	return si, err
}

func (*ImageBuilderClientStub) Ready(ctx context.Context) error {
	return nil
}

func (mock *ImageBuilderClientStub) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	return "ami-0c830793775595d4b-test", nil
}

func (mock *ImageBuilderClientStub) GetAzureImageID(ctx context.Context, composeID string) (string, error) {
	return "/resourceGroups/redhat-deployed/providers/Microsoft.Compute/images/composer-api-92ea98f8-7697-472e-80b1-7454fa0e7fa7", nil
}

func (mock *ImageBuilderClientStub) GetGCPImageName(ctx context.Context, composeID string) (string, error) {
	return "projects/red-hat-image-builder/global/images/composer-api-871fa36d-0b5b-4001-8c95-a11f751a4d66-test", nil
}

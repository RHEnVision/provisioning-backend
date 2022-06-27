package payloads

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/go-chi/render"
)

type InstanceTypeResponse struct {
	*types.InstanceTypeInfo
}

func (s *InstanceTypeResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *InstanceTypeResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListInstanceTypeResponse(sl *[]types.InstanceTypeInfo) []render.Renderer {
	sList := *sl
	list := make([]render.Renderer, 0, len(sList))
	for i := range sList {
		list = append(list, &InstanceTypeResponse{InstanceTypeInfo: &sList[i]})
	}
	return list
}

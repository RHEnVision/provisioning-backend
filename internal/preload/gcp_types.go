package preload

import "fmt"

var GCPInstanceType instanceType

func init() {
	GCPInstanceType = instanceType{
		filename: "gcp_types.yaml",
		path:     "gcp_availability",
		etagName: "gcp-types",
	}
	err := GCPInstanceType.Load()
	if err != nil {
		panic(fmt.Errorf("cannot preload gcp types: %w", err))
	}
}

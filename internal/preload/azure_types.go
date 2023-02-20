package preload

import "fmt"

var AzureInstanceType instanceType

func init() {
	AzureInstanceType = instanceType{
		filename: "azure_types.yaml",
		path:     "azure_availability",
		etagName: "azure-types",
	}
	err := AzureInstanceType.Load()
	if err != nil {
		panic(fmt.Errorf("cannot preload azure types: %w", err))
	}
}

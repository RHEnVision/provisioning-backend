package preload

var GCPInstanceType instanceType

func init() {
	GCPInstanceType = instanceType{
		filename: "gcp_types.yaml",
		path:     "gcp_availability",
		etagName: "gcp-types",
	}
	err := GCPInstanceType.Load()
	if err != nil {
		panic(err)
	}
}

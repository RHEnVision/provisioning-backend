package preload

import "fmt"

var EC2InstanceType instanceType

func init() {
	EC2InstanceType = instanceType{
		filename: "ec2_types.yaml",
		path:     "ec2_availability",
		etagName: "ec2-types",
	}
	err := EC2InstanceType.Load()
	if err != nil {
		panic(fmt.Errorf("cannot preload ec2 types: %w", err))
	}
}

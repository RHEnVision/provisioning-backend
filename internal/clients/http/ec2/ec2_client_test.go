package ec2

import (
	"testing"
)

func TestGetRoleName(t *testing.T) {
	arn := "arn:aws:iam::123456789990:role/role-name"
	roleName, err := getRoleName(arn)
	if err != nil {
		t.Errorf(`rolename parsing error: "%s"`, err)
	} else if roleName != "role-name" {
		t.Error("rolename parsed incorrectly")
	}

	arn = "arn:aws:iam::123456789990:role/service-role/service-role-name"
	_, err = getRoleName(arn)
	if err == nil {
		t.Error("incorrect arn for role parsed as correct")
	}

	arn = "arn:aws:iam::123456789990:role/aws-service-role/test-service.com/aws-service-role-name"
	_, err = getRoleName(arn)
	if err == nil {
		t.Error("incorrect arn for role parsed as correct")
	}

	arn = "arn:aws:iam::123456789990:group/group-name"
	_, err = getRoleName(arn)
	if err == nil {
		t.Error("error expected: group ARN parsed as role")
	}
}

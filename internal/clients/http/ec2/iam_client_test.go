package ec2

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"

	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = Statement{
	Effect: "Allow",
	Action: []string{
		"iam:GetPolicyVersion",
		"iam:GetPolicy",
		"ec2:CreateKeyPair",
		"ec2:CreateLaunchTemplate",
	},
}

var statementDeny = Statement{
	Effect: "Deny",
	Action: []string{
		"iam:GetPolicyVersion",
		"iam:GetPolicy",
		"ec2:CreateKeyPair",
		"ec2:CreateLaunchTemplate",
	},
}

var missingLastPermission = []string{
	"iam:GetPolicyVersion",
	"iam:GetPolicy",
	"ec2:CreateKeyPair",
}

var duplicityPermissions = []string{
	"iam:GetPolicyVersion",
	"iam:GetPolicy",
	"iam:GetPolicyVersion",
	"iam:GetPolicy",
	"ec2:CreateKeyPair",
	"ec2:CreateLaunchTemplate",
}

var actionString = map[string]interface{}{
	"Effect":   "Allow",
	"Resource": "*",
	"Action":   "ec2:StartInstances",
}

var actionArray = map[string]interface{}{
	"Effect":   "Allow",
	"Resource": "*",
	"Action":   []interface{}{"ec2:StartInstances", "ec2:StopInstances"},
}

var noActionArray = map[string]interface{}{
	"Effect":    "Allow",
	"Resource":  "*",
	"NotAction": []interface{}{"ec2:StartInstances", "ec2:StopInstances"},
}

var jsonDataArray = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "iam:GetPolicyVersion",
                "logs:GetLogRecord",
                "iam:ListRolePolicies"
            ],
            "Resource": "*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": "logs:*",
            "Resource": "arn:aws:logs:*:399777895069:log-group:*"
        },
        {
            "Sid": "VisualEditor2",
            "Effect": "Allow",
            "Action": "logs:*",
            "Resource": [
                "arn:aws:logs:*:399777895069:destination:*",
                "arn:aws:logs:*:399777895069:log-group:*:log-stream:*"
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "NotAction": "ec2:*",
            "Resource": "arn:aws:logs:*:399777895069:log-group:*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Deny",
            "Action": "iam:*",
            "Resource": "arn:aws:logs:*:399777895069:log-group:*"
        }
    ]
}`

var jsonData = `{
    "Version": "2012-10-17",
    "Statement": {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "iam:GetPolicyVersion",
                "logs:GetLogRecord",
                "iam:ListRolePolicies"
            ],
            "Resource": "*"
        }
}`

var jsonDataDeny = `{
    "Version": "2012-10-17",
    "Statement": {
            "Sid": "VisualEditor0",
            "Effect": "Deny",
            "Action": [
                "iam:GetPolicyVersion",
                "logs:GetLogRecord",
                "iam:ListRolePolicies"
            ],
            "Resource": "*"
        }
}`

func TestListMissingPermissions(t *testing.T) {
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)

	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	require.NoError(t, err)

	t.Run("get role name", func(t *testing.T) {
		arn := "arn:aws:iam::123456789990:role/role-name"
		roleName, err := getRoleName(arn)
		require.NoError(t, err)
		assert.Equal(t, "role-name", roleName)

		arn = "arn:aws:iam::123456789990:role/service-role/service-role-name"
		_, err = getRoleName(arn)
		require.Error(t, err)

		arn = "arn:aws:iam::123456789990:role/aws-service-role/test-service.com/aws-service-role-name"
		_, err = getRoleName(arn)
		require.Error(t, err)

		arn = "arn:aws:iam::123456789990:group/group-name"
		_, err = getRoleName(arn)
		require.Error(t, err)
	})

	t.Run("list missing permissions", func(t *testing.T) {
		missingStatements := listMissingPermissions(missingLastPermission, expected)
		assert.Equal(t, 1, len(missingStatements))
		assert.Equal(t, missingStatements[0], expected.Action[len(expected.Action)-1])

		shouldBeEmpty := listMissingPermissions(expected.Action, expected)
		assert.Equal(t, 0, len(shouldBeEmpty))

		shouldBeEmpty = listMissingPermissions(duplicityPermissions, expected)
		assert.Equal(t, 0, len(shouldBeEmpty))
	})

	t.Run("get permission from statement", func(t *testing.T) {
		action := getPermissionsFromStatement(ctx, actionString)
		assert.Equal(t, 1, len(action))
		assert.Equal(t, actionString["Action"], action[0])

		actions := getPermissionsFromStatement(ctx, actionArray)
		assert.Equal(t, 2, len(actions))
		assert.Equal(t, "ec2:StartInstances", actions[0])
		assert.Equal(t, "ec2:StopInstances", actions[1])

		noAction := getPermissionsFromStatement(ctx, noActionArray)
		assert.Equal(t, 0, len(noAction))
	})

	t.Run("get permission from allowed", func(t *testing.T) {
		permissions := getPermissionsFromAllowed(ctx, statementDeny)
		assert.Equal(t, 0, len(permissions))
	})

	t.Run("get statement from json", func(t *testing.T) {
		statementArrayPolicies, err := getStatementFromJson(ctx, jsonDataArray)
		require.NoError(t, err)
		assert.Equal(t, 5, len(statementArrayPolicies))

		statementPolicies, err := getStatementFromJson(ctx, jsonData)
		require.NoError(t, err)
		assert.Equal(t, 3, len(statementPolicies))

		statementDeny, err := getStatementFromJson(ctx, jsonDataDeny)
		require.NoError(t, err)
		assert.Equal(t, 0, len(statementDeny))
	})
}

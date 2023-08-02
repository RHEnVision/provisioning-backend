package ec2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

var AWSPermissionsMissing = errors.New("AWS permissions missing")

// Statement is a main policy element.
type Statement struct {
	// Sid is an optional identifier of the Statement
	Sid string

	// Effect defines results outcome
	// Values: "Allow" or "Deny"
	Effect string

	// Describes the specific AWS actions for which the effect applies.
	Action []string
}

// Policy AWS JSON element
type Policy struct {
	Statement []Statement `json:"Statement"`
}

var expectedStatement = Statement{
	Sid:    "RedHatProvisioning",
	Effect: "Allow",
	Action: []string{
		"iam:GetPolicyVersion",
		"iam:GetPolicy",
		"iam:ListAttachedRolePolicies",
		"iam:GetRolePolicy",
		"ec2:CreateKeyPair",
		"ec2:CreateLaunchTemplate",
		"ec2:CreateLaunchTemplateVersion",
		"ec2:CreateTags",
		"ec2:DeleteKeyPair",
		"ec2:DeleteTags",
		"ec2:DescribeAvailabilityZones",
		"ec2:DescribeImages",
		"ec2:DescribeInstanceTypes",
		"ec2:DescribeInstances",
		"ec2:DescribeKeyPairs",
		"ec2:DescribeLaunchTemplates",
		"ec2:DescribeLaunchTemplateVersions",
		"ec2:DescribeRegions",
		"ec2:DescribeSecurityGroups",
		"ec2:DescribeSnapshotAttribute",
		"ec2:DescribeTags",
		"ec2:ImportKeyPair",
		"ec2:RunInstances",
		"ec2:StartInstances",
		"iam:ListRolePolicies",
	},
}

func getRoleName(arn string) (string, error) {
	arnParts := strings.Split(arn, ":")
	if len(arnParts) == 0 {
		return "", fmt.Errorf("%w: ARN has no colons: %s", http.ARNParsingError, arn)
	}
	roleName := strings.Split(arnParts[len(arnParts)-1], "/")
	if len(roleName) != 2 {
		return "", fmt.Errorf("%w: ARN has incorrect syntax: %s rolename parsing result: %s",
			http.ARNParsingError, arn, roleName)
	} else if roleName[0] != "role" {
		return "", fmt.Errorf("%w: ARN does not have any role: %s", http.ARNParsingError, arn)
	}
	return roleName[1], nil
}

func getPermissionsFromAllowed(ctx context.Context, statementBody interface{}) []string {
	switch statementFields := statementBody.(type) {
	case map[string]interface{}:
		switch effectField := statementFields["Effect"].(type) {
		case string:
			if effectField == "Allow" {
				return getPermissionsFromStatement(ctx, statementFields)
			}
		}
	}
	return nil
}

func getPermissionsFromStatement(ctx context.Context, statementBody map[string]interface{}) []string {
	var result []string

	switch actionField := statementBody["Action"].(type) {
	case string:
		result = append(result, actionField)
	case []interface{}:
		for _, v := range actionField {
			switch actionElement := v.(type) {
			case string:
				result = append(result, actionElement)
			}
		}
	}
	return result
}

func getJsonFromAWSDocument(ctx context.Context, document string) (string, error) {
	unescapedDoc, err := url.QueryUnescape(document)
	if err != nil {
		return "", fmt.Errorf("could not parse URL from the document: %w", err)
	}
	return unescapedDoc, nil
}

func getStatementFromJson(ctx context.Context, document string) ([]string, error) {
	var result []string

	jsonBody := []byte(document)
	var policy interface{}
	err := json.Unmarshal(jsonBody, &policy)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %w", err)
	}

	switch policyBody := policy.(type) {
	case map[string]interface{}:
		switch statementField := policyBody["Statement"].(type) {
		case []interface{}:
			for _, v := range statementField {
				switch statementElement := v.(type) {
				case interface{}:
					result = append(result, getPermissionsFromAllowed(ctx, statementElement)...)
				}
			}
		case interface{}:
			result = append(result, getPermissionsFromAllowed(ctx, statementField)...)

		}
	}
	return result, nil
}

func listStatements(ctx context.Context, versions []*iamTypes.PolicyVersion) ([]string, error) {
	var result []string
	for _, version := range versions {
		jsonDocument, err := getJsonFromAWSDocument(ctx, *version.Document)
		if err != nil {
			return nil, fmt.Errorf("could not get JSON from AWS document: %w", err)
		}

		policyVersionStatements, err := getStatementFromJson(ctx, jsonDocument)
		if err != nil {
			return nil, fmt.Errorf("could not fetch statement from AWS document: %w", err)
		}

		result = append(result, policyVersionStatements...)
	}
	return result, nil
}

func listMissingPermissions(statements []string, expected Statement) []string {
	presentPermissions := make(map[string]struct{})
	var missing []string
	for _, statement := range statements {
		presentPermissions[statement] = struct{}{}
	}

	for _, statement := range expected.Action {
		if _, ok := presentPermissions[statement]; !ok {
			missing = append(missing, statement)
		}
	}
	return missing
}

func (c *ec2Client) listAttachedRolePolicies(ctx context.Context, roleName string) ([]*iamTypes.AttachedPolicy, error) {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	output, err := c.iam.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("cannot list attached role policies: %w", err)
	}

	result := make([]*iamTypes.AttachedPolicy, len(output.AttachedPolicies))
	for i := range output.AttachedPolicies {
		result[i] = &output.AttachedPolicies[i]
	}
	return result, nil
}

func (c *ec2Client) listInlineRolePolicies(ctx context.Context, roleName string) ([]string, error) {
	rolePoliciesInput := &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	output, err := c.iam.ListRolePolicies(ctx, rolePoliciesInput)
	if err != nil {
		return nil, fmt.Errorf("cannot list inline role policies: %w", err)
	}

	result := make([]string, len(output.PolicyNames))
	for i := range output.PolicyNames {
		getRolePolicyInput := &iam.GetRolePolicyInput{
			RoleName:   aws.String(roleName),
			PolicyName: &output.PolicyNames[i],
		}
		getRolePolicyOutput, err := c.iam.GetRolePolicy(ctx, getRolePolicyInput)
		if err != nil {
			return nil, fmt.Errorf("cannot get inline role policy: %w", err)
		}
		result[i] = *getRolePolicyOutput.PolicyDocument

	}
	return result, nil
}

func (c *ec2Client) listPoliciesFromAttached(ctx context.Context, policies []*iamTypes.AttachedPolicy) ([]*iamTypes.Policy, error) {
	result := make([]*iamTypes.Policy, len(policies))
	for i := range policies {
		input := &iam.GetPolicyInput{
			PolicyArn: policies[i].PolicyArn,
		}
		output, err := c.iam.GetPolicy(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("cannot get policy for policy arn %s: %w", *input.PolicyArn, err)
		}
		result[i] = output.Policy

	}
	return result, nil
}

func (c *ec2Client) listPolicyVersions(ctx context.Context, policies []*iamTypes.Policy) ([]*iamTypes.PolicyVersion, error) {
	result := make([]*iamTypes.PolicyVersion, len(policies))
	for i := range policies {
		input := &iam.GetPolicyVersionInput{
			PolicyArn: policies[i].Arn,
			VersionId: policies[i].DefaultVersionId,
		}
		output, err := c.iam.GetPolicyVersion(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("cannot get policy version for policy arn %s, version %s: %w", *input.PolicyArn, *input.VersionId, err)
		}
		result[i] = output.PolicyVersion
	}
	return result, nil
}

func (c *ec2Client) checkInlinePolicies(ctx context.Context, missingPermissions []string, roleName string) ([]string, error) {
	inlinePoliciesStatement, err := c.listInlineRolePolicies(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not list inline policy documents: %w", err)
	}
	missingStatement := Statement{
		Effect: "Allow",
		Action: missingPermissions,
	}
	missing := listMissingPermissions(inlinePoliciesStatement, missingStatement)
	if len(missingPermissions) != 0 {
		return missing, fmt.Errorf("%w: %s", AWSPermissionsMissing, strings.Join(missing, ", "))
	}
	return nil, nil
}

func (c *ec2Client) CheckPermission(ctx context.Context, auth *clients.Authentication) ([]string, error) {
	roleName, err := getRoleName(auth.Payload)
	if err != nil {
		return nil, fmt.Errorf("unable to parse ARN: %w", err)
	}

	attachedRolePolicies, err := c.listAttachedRolePolicies(ctx, roleName)
	if err != nil {
		return nil, err
	}
	policies, err := c.listPoliciesFromAttached(ctx, attachedRolePolicies)
	if err != nil {
		return nil, fmt.Errorf("could not get policies: %w", err)
	}

	versions, err := c.listPolicyVersions(ctx, policies)
	if err != nil {
		return nil, fmt.Errorf("could not get policy versions: %w", err)
	}

	statements, err := listStatements(ctx, versions)
	if err != nil {
		return nil, fmt.Errorf("could not list statements: %w", err)
	}

	missingPermissions := listMissingPermissions(statements, expectedStatement)
	if len(missingPermissions) != 0 {
		return c.checkInlinePolicies(ctx, missingPermissions, roleName)
	}

	return nil, nil
}

package clients

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed http/azure/lighthouse.tmpl.json
var offeringTemplate string

type AzureOfferingTemplate struct {
	// OfferingDefaultName that Customer can change while deploying the offering
	OfferingDefaultName string

	// OfferingDefaultDescription describing the offering, can be changed by Customer while deploying
	OfferingDefaultDescription string

	// TenantID of the offering tenant (Azure account)
	TenantID string

	// EnterpriseAppID of the App that will act as an offering Principal
	EnterpriseAppID string

	// EnterpriseAppName of the offering principal - the display name
	EnterpriseAppName string
}

func (tempParams AzureOfferingTemplate) Render(ctx context.Context, wr io.Writer) error {
	tmpl, err := template.New("lighthouse.tmpl.json").Parse(offeringTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse Azure template: %w", err)
	}
	if err = tmpl.Execute(wr, tempParams); err != nil {
		return fmt.Errorf("failed to render Azure offering: %w", err)
	}
	return nil
}

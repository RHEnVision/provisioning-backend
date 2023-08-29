// Package init provides functions for integration with virtual machine initialization
// frameworks like cloud-init.
package userdata

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/rs/zerolog"
)

type UserData struct {
	// Type defines hyperscaler platform for which the user data should be generated for
	Type models.ProviderType

	// PowerOff controls power state first boot setting. When set, user data
	// will contain instruction to power off the VM after initial launch after
	// initialization.
	PowerOff bool

	// PowerOffDelayMin specifies number of minutes after shutdown command is
	// executed. A warning over tty is sent to all logged users. Minimum value
	// is 1, when set negative or to zero, one minute will be used.
	PowerOffDelayMin int

	// A message that is passed to the poweroff command and shown to logged users.
	// When unset, "User data scheduled power off" is used.
	PowerOffMessage string

	// InsightsTags renders a first-boot script which populates /etc/insights-client/tags.yaml
	InsightsTags bool
}

func (ud UserData) IsAWS() bool {
	return ud.Type == models.ProviderTypeAWS
}

func (ud UserData) IsAzure() bool {
	return ud.Type == models.ProviderTypeAzure
}

func (ud UserData) IsGCP() bool {
	return ud.Type == models.ProviderTypeGCP
}

//go:embed cloud-init.goyaml
var cloudinitBuffer []byte
var cloudinitTemplate *template.Template

//go:embed script.tmpl
var scriptBuffer []byte
var scriptTemplate *template.Template

func init() {
	var err error
	cloudinitTemplate, err = template.New("cloudinit").Parse(string(cloudinitBuffer))
	if err != nil {
		panic(err)
	}
	scriptTemplate, err = template.New("script").Parse(string(scriptBuffer))
	if err != nil {
		panic(err)
	}
}

// GenerateUserData creates a cloud-init user-data from a build-in template.
func GenerateUserData(ctx context.Context, userData *UserData) ([]byte, error) {
	logger := zerolog.Ctx(ctx)

	if userData.PowerOffDelayMin < 1 {
		userData.PowerOffDelayMin = 1
	}
	if userData.PowerOffMessage == "" {
		userData.PowerOffMessage = "User data scheduled power off"
	}

	var buffer bytes.Buffer
	var err error
	if userData.Type == models.ProviderTypeGCP {
		err = scriptTemplate.Execute(&buffer, userData)
	} else {
		err = cloudinitTemplate.Execute(&buffer, userData)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot generate user data: %w", err)
	}

	udBytes := buffer.Bytes()
	logger.Trace().Bytes("payload", udBytes).Msg("Generated userdata")
	return udBytes, nil
}

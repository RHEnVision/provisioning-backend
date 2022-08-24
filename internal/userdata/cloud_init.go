// Package init provides functions for integration with virtual machine initialization
// frameworks like cloud-init.
package userdata

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

type UserData struct {
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
}

//go:embed userdata.goyaml
var userdataBuffer []byte
var userdataTemplate *template.Template

func init() {
	var err error
	userdataTemplate, err = template.New("userdata").Parse(string(userdataBuffer))
	if err != nil {
		panic(err)
	}
}

// GenerateUserData creates a cloud-init user-data from a build-in template.
func GenerateUserData(userData *UserData) ([]byte, error) {
	if userData.PowerOffDelayMin < 1 {
		userData.PowerOffDelayMin = 1
	}
	if userData.PowerOffMessage == "" {
		userData.PowerOffMessage = "User data scheduled power off"
	}
	var buffer bytes.Buffer
	err := userdataTemplate.Execute(&buffer, userData)
	if err != nil {
		return nil, fmt.Errorf("cannot generate user data: %w", err)
	}
	return buffer.Bytes(), nil
}

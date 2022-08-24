package userdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDefaults(t *testing.T) {
	userDataInput := UserData{}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config
`
	assert.Equal(t, expected, string(userData))
}

func TestGeneratePoweroff(t *testing.T) {
	userDataInput := UserData{
		PowerOff: true,
	}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config
power_state:
  mode: poweroff
  delay: "+1"
  message: "User data scheduled power off"
  timeout: 60
`
	assert.Equal(t, expected, string(userData))
}

func TestGenerateCustomPoweroff(t *testing.T) {
	userDataInput := UserData{
		PowerOff:         true,
		PowerOffDelayMin: 42,
		PowerOffMessage:  "Blah",
	}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config
power_state:
  mode: poweroff
  delay: "+42"
  message: "Blah"
  timeout: 60
`
	assert.Equal(t, expected, string(userData))
}

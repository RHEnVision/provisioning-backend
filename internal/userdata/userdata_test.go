package userdata

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var trimRe, _ = regexp.Compile(`\n{2,}`)

func validateYAML(str []byte) error {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(str, &m)
	if err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	return nil
}

func TestGenerateDefaults(t *testing.T) {
	userDataInput := UserData{}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config`

	assert.NoError(t, validateYAML(userData))
	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
}

func TestGenerateAWSTags(t *testing.T) {
	userDataInput := UserData{
		Type:         models.ProviderTypeAWS,
		InsightsTags: true,
	}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config
write_files:
- path: /etc/insights-client/tags-generate.sh
  owner: root:root
  permissions: '0770'
  content: |
    #!/bin/sh
    TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
    PUBLIC_IP4=$(/usr/bin/curl -sH "X-aws-ec2-metadata-token: $TOKEN" --connect-timeout 5 http://169.254.169.254/latest/meta-data/public-ipv4)
    PUBLIC_HOSTNAME=$(/usr/bin/curl -sH "X-aws-ec2-metadata-token: $TOKEN" --connect-timeout 5 http://169.254.169.254/latest/meta-data/public-hostname)
    test -d /etc/insights-client || mkdir /etc/insights-client
    echo "---" > /etc/insights-client/tags.yaml
    echo "Public hostname: $PUBLIC_HOSTNAME" >> /etc/insights-client/tags.yaml
    echo "Public IPv4: $PUBLIC_IP4" >> /etc/insights-client/tags.yaml
runcmd:
- [ "/bin/sh", "-xc", "/etc/insights-client/tags-generate.sh" ]`

	assert.NoError(t, validateYAML(userData))
	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
}

func TestGenerateAzureTags(t *testing.T) {
	userDataInput := UserData{
		Type:         models.ProviderTypeAzure,
		InsightsTags: true,
	}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#cloud-config
write_files:
- path: /etc/insights-client/tags-generate.sh
  owner: root:root
  permissions: '0770'
  content: |
    #!/bin/sh
    PUBLIC_IP4=$(curl -s -H Metadata:true --noproxy "*" "http://169.254.169.254/metadata/instance?api-version=2021-02-01" | /usr/libexec/platform-python -c 'import json,sys;print(json.load(sys.stdin)["network"]["interface"][0]["ipv4"]["ipAddress"][0]["publicIpAddress"])')
    LOADBALANCER_IP4=$(/usr/bin/curl -sH "Metadata:true" --connect-timeout 5 http://169.254.169.254/metadata/loadbalancer?api-version=2020-10-01 | /usr/libexec/platform-python -c 'import json,sys;print(json.load(sys.stdin)["loadbalancer"]["publicIpAddresses"][0]["frontendIpAddress"])' 2>/dev/null)
    test -d /etc/insights-client || mkdir /etc/insights-client
    echo "---" > /etc/insights-client/tags.yaml
    echo "Public IPv4: $PUBLIC_IP4" >> /etc/insights-client/tags.yaml
    echo "Public LB IPv4: $LOADBALANCER_IP4" >> /etc/insights-client/tags.yaml
runcmd:
- [ "/bin/sh", "-xc", "/etc/insights-client/tags-generate.sh" ]`

	assert.NoError(t, validateYAML(userData))
	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
}

func TestGenerateGCPTags(t *testing.T) {
	userDataInput := UserData{
		Type:         models.ProviderTypeGCP,
		InsightsTags: true,
	}
	userData, err := GenerateUserData(&userDataInput)
	require.NoError(t, err)
	expected := `#! /bin/bash
PUBLIC_IP4=$(/usr/bin/curl -sH "Metadata-Flavor: Google" --connect-timeout 5 http://metadata/computeMetadata/v1/instance/network-interfaces/0/ip)
test -d /etc/insights-client || mkdir /etc/insights-client
echo "---" > /etc/insights-client/tags.yaml
echo "Public IPv4: $PUBLIC_IP4" >> /etc/insights-client/tags.yaml
exit 0`

	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
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
  timeout: 60`

	assert.NoError(t, validateYAML(userData))
	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
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
  timeout: 60`

	assert.NoError(t, validateYAML(userData))
	assert.Equal(t, expected, strings.Trim(trimRe.ReplaceAllString(string(userData), "\n"), "\n"))
}

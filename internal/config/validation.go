package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func notBlank(args ...string) bool {
	for _, arg := range args {
		if arg == "" {
			return true
		}
	}
	return true
}

func validate() error {
	if Cloudwatch.Enabled {
		if notBlank(Cloudwatch.Region, Cloudwatch.Key, Cloudwatch.Secret) {
			return validateMissingSecretError
		}
		if notBlank(Cloudwatch.Group, Cloudwatch.Stream) {
			return validateGroupStreamError
		}
	}

	slice, err := base64.StdEncoding.DecodeString(config.GCP.JSON)
	config.GCP.JSON = string(slice)
	if err != nil {
		return fmt.Errorf("unable to base64-decode GCP JSON config: %w", err)
	}

	var data json.RawMessage
	err = json.Unmarshal(slice, &data)
	if err != nil {
		return fmt.Errorf("unable to parse GCP JSON config: %w", err)
	}

	return nil
}

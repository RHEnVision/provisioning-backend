package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func validate() error {
	if Cloudwatch.Enabled {
		if Cloudwatch.Region == "" || Cloudwatch.Key == "" || Cloudwatch.Secret == "" {
			return validateMissingSecretError
		}
		if Cloudwatch.Group == "" || Cloudwatch.Stream == "" {
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

//go:build test
// +build test

package config

import "github.com/spf13/viper"

func init() {
	viper.Set("featureFlags.environment", "test")
}

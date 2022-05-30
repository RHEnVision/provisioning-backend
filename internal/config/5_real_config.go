//go:build !test
// +build !test

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("featureFlags.environment", "developement")

	viper.AddConfigPath("./configs")
	viper.SetConfigName("defaults")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.SetConfigName("local")
	err = viper.MergeInConfig()
	if err != nil {
		fmt.Println("Could not read local.yaml", err)
	}
}

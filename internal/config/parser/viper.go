// This package provides viper parser.
//
// Due to bug in v1 API, it is not possible to set custom parser in the global viper
// package. For this reason, all parsing must be done in a separate package where
// viper is initialized with custom StringReplacer.
//
// https://github.com/spf13/viper/pull/870
//
package parser

import "github.com/spf13/viper"

var Viper *viper.Viper

func init() {
	Viper = viper.NewWithOptions(viper.EnvKeyReplacer(customReplacer{}))
	Viper.AutomaticEnv()
}

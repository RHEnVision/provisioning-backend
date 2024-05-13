//go:build integration && !test

package config

import (
	"os"
)

func init() {
	os.Setenv("APP_CACHE_TYPE_APP_ID", "false")
	os.Setenv("APP_CACHE_ACCOUNT", "false")
	os.Setenv("FEATURE_FLAGS_ENVIRONMENT", "test")
	os.Setenv("TEST_ENVIRONMENT", "integration")
	os.Setenv("WORKER_QUEUE", "memory")
}

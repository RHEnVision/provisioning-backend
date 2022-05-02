package config

import (
	"os"
	"strconv"
	"sync"
)

type LoggingConfig struct {
	// Global level (-1 trace, 0 debug, 1 info ... 5 panic)
	Level       int
	ExitOnPanic bool
	Cloudwatch  bool
	Stdout      bool
	CWGroup     string
	CWStream    string
	AWSRegion   string
	AWSKey      string
	AWSSecret   string
	AWSSession  string

	initialized bool
}

var loggingConfig LoggingConfig
var initMutex sync.Mutex

func initializeConfig() {
	level, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	loggingConfig = LoggingConfig{
		Level:       level,
		ExitOnPanic: os.Getenv("EXIT_ON_PANIC") == "1",
		Stdout:      os.Getenv("STDOUT_LOG") == "1",
		Cloudwatch:  os.Getenv("CLOUDWATCH_LOG") == "1",
		CWGroup:     os.Getenv("CLOUDWATCH_GROUP"),
		CWStream:    os.Getenv("CLOUDWATCH_STREAM"),
		AWSRegion:   os.Getenv("AWS_REGION"),
		AWSKey:      os.Getenv("AWS_KEY"),
		AWSSecret:   os.Getenv("AWS_SECRET"),
		AWSSession:  os.Getenv("AWS_SESSION"),
		initialized: true,
	}
}

func GetLoggingConfig() *LoggingConfig {
	initMutex.Lock()
	defer initMutex.Unlock()

	if !loggingConfig.initialized {
		initializeConfig()
	}

	return &loggingConfig
}

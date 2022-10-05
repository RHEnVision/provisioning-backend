package ec2

import (
	"context"

	"github.com/aws/smithy-go/logging"
	"github.com/rs/zerolog"
)

type ec2Logger struct {
	zlog *zerolog.Logger
}

func (logger *ec2Logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	logger.zlog.Trace().Msgf(format, v...)
}

func NewEC2Logger(ctx context.Context) *ec2Logger {
	logger := logger(ctx)
	return &ec2Logger{
		zlog: &logger,
	}
}

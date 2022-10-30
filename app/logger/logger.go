package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		// backup
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logger.(*logrus.Entry)
}

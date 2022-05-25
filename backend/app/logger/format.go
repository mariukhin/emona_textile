package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type LogFormat string

const (
	Text LogFormat = "text"
	Json LogFormat = "json"
)

var logFormats = []LogFormat{Text, Json}

func (l LogFormat) Formatter() logrus.Formatter {
	switch l {
	case Text:
		return &logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		}
	case Json:
		return &logrus.JSONFormatter{}

	default:
		panic(fmt.Sprintf("unsupported log format %s", l))
	}
}

func ParseLogFormat(value string) (LogFormat, error) {
	for _, format := range logFormats {
		if format == LogFormat(value) {
			return format, nil
		}
	}

	return "", fmt.Errorf("unknown log format value '%s'", value)
}

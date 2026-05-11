package esl

import "github.com/panjf2000/gnet/v2/pkg/logging"

func isTraceEnabled() bool {
	return logging.DebugLevel >= options.Level
}

func isDebugEnabled() bool {
	return logging.DebugLevel >= options.Level
}

func isInfoEnabled() bool {
	return logging.InfoLevel >= options.Level
}

package esl

import "github.com/bytedance/gopkg/util/logger"

func isTraceEnabled() bool {
	return logger.LevelTrace >= options.Level
}

func isDebugEnabled() bool {
	return logger.LevelDebug >= options.Level
}

func isInfoEnabled() bool {
	return logger.LevelInfo >= options.Level
}

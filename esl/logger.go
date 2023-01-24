package esl

import (
	"fmt"
)

var (
//log = logging.MustGetLogger("goesl")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
// format = logging.MustStringFormatter(
//
//	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s}%{color:reset} %{message}",
//
// )
)

func Debug(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func Error(message string, args ...interface{}) {
	fmt.Errorf(message, args...)
}

func Notice(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func Info(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func Warning(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func init() {
	//backend := logging.NewLogBackend(os.Stderr, "", 0)
	//formatter := logging.NewBackendFormatter(backend, format)
	//logging.SetBackend(formatter)
}

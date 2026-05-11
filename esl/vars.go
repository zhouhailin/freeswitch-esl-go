package esl

import "go.uber.org/zap/zapcore"

type Name string

const (
	NEW_LINE byte = 10 // 或者 0x0A

	MESSAGE_TERMINATOR     = "\n\n"
	LINE_TERMINATOR        = "\n"
	OK                     = "+OK"
	AUTH_REQUEST           = "auth/request"
	API_RESPONSE           = "api/response"
	COMMAND_REPLY          = "command/reply"
	TEXT_EVENT_PLAIN       = "text/event-plain"
	TEXT_EVENT_XML         = "text/event-xml"
	TEXT_DISCONNECT_NOTICE = "text/disconnect-notice"
	TEXT_RUDE_REJECTION    = "text/rude-rejection"
	ERR_INVALID            = "-ERR invalid"
)

// Level is the alias of zapcore.Level.
type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

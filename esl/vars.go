package esl

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

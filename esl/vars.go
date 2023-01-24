package esl

var (

	// Size of buffer when we read from connection.
	// 1024 << 6 == 65536
	ReadBufferSize = 1024 << 6

	// Freeswitch events that we can handle (have logic for it)
	AvailableMessageTypes = []string{"auth/request", "text/disconnect-notice", "text/event-json", "text/event-plain", "api/response", "command/reply"}
)

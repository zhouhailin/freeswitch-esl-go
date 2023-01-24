package esl

var (
	EInvalidCommandProvided  = "Invalid command provided. Command cannot contain \\r and/or \\n. Provided command is: %s"
	ECouldNotReadMIMEHeaders = "Error while reading MIME headers: %s"
	EInvalidContentLength    = "Unable to get size of content-length: %s"
	EUnsuccessfulReply       = "Got error while reading from reply command: %s"
	ECouldNotReadyBody       = "Got error while reading reader body: %s"
	EUnsupportedMessageType  = "Unsupported message type! We got '%s'. Supported types are: %v "
	ECouldNotDecode          = "Could not decode/unescape message: %s"
	ECouldNotStartListener   = "Got error while attempting to start listener: %s"
	EListenerConnection      = "Listener connection error: %s"
	EInvalidServerAddr       = "Please make sure to pass along valid address. You've passed: \"%s\""
	EUnexpectedAuthHeader    = "Expected auth/request content type. Got %s"
	EInvalidPassword         = "Could not authenticate against freeswitch with provided password: %s"
	ECouldNotCreateMessage   = "Error while creating new message: %s"
	ECouldNotSendEvent       = "Must send at least one event header, detected `%d` header"
)

package esl

// IEslEventListener - IEslEventListener
type IEslEventListener interface {
	// EventReceived - Signal of a server initiated event.
	EventReceived(event *EslEvent) error
	// BackgroundJobResultReceived - Signal of an event containing the result of a client requested background job.
	BackgroundJobResultReceived(event *EslEvent) error
}

// IEslProtocolListener - IEslProtocolListener
type IEslProtocolListener interface {
	// Signal of a server initiated event.
	authResponseReceived(c *Client, response *CommandResponse)
	// Signal of an event containing the result of a client requested background job.
	eventReceived(c *Client, event *EslEvent)
	// disconnected.
	disconnected(c *Client)
}

// IEslConnectionListener - Esl Connection Listener
type IEslConnectionListener interface {

	// ConnectFailure - Connect failure
	ConnectFailure(c *Client)

	// Connected - success
	Connected(c *Client)

	// Authenticated - authentication
	Authenticated(authenticated bool, c *Client)

	// Disconnected - connection is closed
	Disconnected(c *Client)
}

package esl

import "strings"

// CommandResponse server's response
type CommandResponse struct {
	command   string
	replyText string
	response  *EslMessage
	success   bool
}

func NewCommandResponse(command string, response *EslMessage) *CommandResponse {
	replyText := response.GetHeaderValue(REPLY_TEXT)
	return &CommandResponse{
		command:   command,
		replyText: replyText,
		response:  response,
		success:   strings.HasPrefix(replyText, OK),
	}
}

// GetCommand - original command sent to the server
// @return the original command sent to the server
func (resp *CommandResponse) GetCommand() string {
	return resp.command
}

// IsOk - the response Reply-Text line starts with "+OK"
// @return true if and only if the response Reply-Text line starts with "+OK"
func (resp *CommandResponse) IsOk() bool {
	return resp.success
}

// GetReplyText - the response Reply-Text line
// @return the full response Reply-Text line.
func (resp *CommandResponse) GetReplyText() string {
	return resp.replyText
}

// GetResponse - the full response from the server
// @return {@link EslMessage} the full response from the server.
func (resp *CommandResponse) GetResponse() *EslMessage {
	return resp.response
}

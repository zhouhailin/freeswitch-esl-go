package esl

import "strings"

// SendEvent Event
type SendEvent struct {
	msgLines []string
}

// NewSendEvent - Constructor for use with outbound socket client only.  This client mode does not need a call UUID for context.
// @param name part of line
func NewSendEvent(name string) *SendEvent {
	e := SendEvent{msgLines: *new([]string)}
	e.msgLines = append(e.msgLines, "sendevent "+name)
	return &e
}

// AddLine - A generic method to add a message line. The constructed line in the sent message will be in the
//   - form:
//   - <pre>
//   - name: value
//   - </pre>
//   - @param name  part of line
//   - @param value part of line
//   - @return a {@link SendEvent} object.
func (e *SendEvent) AddLine(name, value string) {
	e.msgLines = append(e.msgLines, name+": "+value)
}

// AddBody - A generic method to add a message line. The constructed line in the sent message will be in the
//   - form:
//   - <pre>
//   - name: value
//   - </pre>
//   - @param line part of line
func (e *SendEvent) AddBody(line string) {
	e.msgLines = append(e.msgLines, line)
}

// GetMsgLines - The list of strings that make up the message to send to NextSWITCH.
//   - @return list of strings, as they were added to this message.
func (e *SendEvent) GetMsgLines() *[]string {
	return &e.msgLines
}

// ToString - To String
func (e *SendEvent) ToString() string {
	var sb strings.Builder
	sb.WriteString("SendEvent: ")
	if len(e.msgLines) > 0 {
		sb.WriteString(e.msgLines[0])
	}
	return sb.String()
}

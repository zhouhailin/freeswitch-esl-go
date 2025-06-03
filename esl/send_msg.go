package esl

import (
	"strconv"
)

type SendMsg struct {
	lines   []string
	hasUuid bool
}

// NewSendMsg - Constructor for use with the inbound client.
func NewSendMsg(uuid string) *SendMsg {
	sendMsg := SendMsg{lines: *new([]string)}
	if uuid == "" {
		sendMsg.lines = append(sendMsg.lines, "sendmsg")
		sendMsg.hasUuid = false
	} else {
		sendMsg.lines = append(sendMsg.lines, "sendmsg "+uuid)
		sendMsg.hasUuid = true
	}
	return &sendMsg
}

// AddCallCommand - Adds the following line to the message:
//   - <pre>
//   - call-command: command
//   - </pre>
//   - @param command the string command [ execute | hangup ]
func (m *SendMsg) AddCallCommand(command string) {
	m.lines = append(m.lines, "call-command: "+command)
}

// AddExecuteAppName -  Adds the following line to the message:
//   - <pre>
//   - execute-app-name: appName
//   - </pre>
//   - @param appName the string app name to execute
func (m *SendMsg) AddExecuteAppName(appName string) {
	m.lines = append(m.lines, "execute-app-name: "+appName)
}

// AddExecuteAppArg -  Adds the following line to the message:
//   - <pre>
//   - execute-app-arg: arg
//   - </pre>
//   - @param arg the string arg
func (m *SendMsg) AddExecuteAppArg(arg string) {
	m.lines = append(m.lines, "execute-app-arg: "+arg)
}

// AddLoops - Adds the following line to the message:
//   - <pre>
//   - loops: count
//   - </pre>
//   - @param count the int number of times to loop
func (m *SendMsg) AddLoops(count int) {
	m.lines = append(m.lines, "loops: "+strconv.Itoa(count))
}

// AddHangupCause - Adds the following line to the message:
//   - <pre>
//   - hangup-cause: cause
//   - </pre>
//   - @param cause the string cause
func (m *SendMsg) AddHangupCause(cause string) {
	m.lines = append(m.lines, "hangup-cause: "+cause)
}

// AddNomediaUuid - Adds the following line to the message:
//   - <pre>
//   - nomedia-uid: value
//   - </pre>
//   - @param value the string value part of the line
func (m *SendMsg) AddNomediaUuid(value string) {
	m.lines = append(m.lines, "nomedia-uuid: "+value)
}

// AddEventLock - Adds the following line to the message:
//   - <pre>
//   - event-lock: true
//   - </pre>
func (m *SendMsg) AddEventLock() {
	m.lines = append(m.lines, "event-lock: true")
}

// AddGenericLine - A generic method to add a message line. The constructed line in the sent message will be in the
//   - form:
//   - <pre>
//   - name: value
//   - </pre>
//   - @param name part of line
//   - @param value part of line
func (m *SendMsg) AddGenericLine(name, value string) {
	m.lines = append(m.lines, name+": "+value)
}

// GetMsgLines - The list of strings that make up the message to send to FreeSWITCH.
//   - @return list of strings, as they were added to this message.
func (m *SendMsg) GetMsgLines() *[]string {
	return &m.lines
}

// HasUuid - Indicate if message was constructed with a UUID.
//   - @return true if constructed with a UUID.
func (m *SendMsg) HasUuid() bool {
	return m.hasUuid
}

// ToString - To String.
func (m *SendMsg) ToString() string {
	if len(m.lines) > 1 {
		return m.lines[1]
	}
	return m.lines[0]
}

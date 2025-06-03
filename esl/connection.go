package esl

import (
	"errors"
	"github.com/cloudwego/netpoll"
	"net"
	"strings"
	"sync"
)

// SocketConnection Main connection against ESL - Gotta add more description here
type SocketConnection struct {
	netpoll.Connection
	msg                    chan *EslMessage
	mtx                    sync.Mutex
	eventListeners         []IEslEventListener
	connectionListeners    []IEslConnectionListener
	authenticationResponse *CommandResponse
	authenticatorResponded bool
	authenticated          bool
	listener               IEslProtocolListener
}

func (socket *SocketConnection) CanSend() bool {
	return socket != nil && socket.Connection != nil && socket.IsActive() && socket.authenticated
}

func (socket *SocketConnection) AddEventListener(listener IEslEventListener) {
	if socket.eventListeners == nil {
		socket.eventListeners = *new([]IEslEventListener)
	}
	socket.eventListeners = append(socket.eventListeners, listener)
}

func (socket *SocketConnection) AddConnectionListener(listener IEslConnectionListener) {
	if socket.connectionListeners == nil {
		socket.connectionListeners = *new([]IEslConnectionListener)
	}
	socket.connectionListeners = append(socket.connectionListeners, listener)
}

// SendSyncApiCommand Sends a NextSWITCH API command to the server and blocks, waiting for an immediate response from the server.
func (socket *SocketConnection) SendSyncApiCommand(command, arg string) (*EslMessage, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if command != "" {
		sb.WriteString("bgapi ")
		sb.WriteString(command)
	}
	if arg != "" {
		sb.WriteString(" ")
		sb.WriteString(arg)
	}
	return socket.sendSyncSingleLineCommand(sb.String())
}

// SendAsyncApiCommand Submit a NextSWITCH API command to the server to be executed in background mode.
// A synchronous response from the server provides a UUID to identify the job execution results.
// When the server has completed the job execution it fires a BACKGROUND_JOB Event with the execution results.
func (socket *SocketConnection) SendAsyncApiCommand(command, arg string) (*string, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if command != "" {
		sb.WriteString("bgapi ")
		sb.WriteString(command)
	}
	if arg != "" {
		sb.WriteString(" ")
		sb.WriteString(arg)
	}
	return socket.sendAsyncCommand(sb.String())
}

// SetEventSubscriptions Set the current event subscription for this connection to the server.
// Examples of the events argument are:
//
//	ALL
//	CHANNEL_CREATE CHANNEL_DESTROY HEARTBEAT
//	CHANNEL_CREATE CHANNEL_DESTROY CUSTOM conference::maintenance sofia::register sofia::expire
//
// format - can be { plain | xml }
// events { all | space separated list of events }
func (socket *SocketConnection) SetEventSubscriptions(format, events string) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	if format != "plain" {
		return nil, errors.New("Only 'plain' event format is supported at present")
	}
	command := "event " + format + " " + events
	response, err := socket.sendSyncSingleLineCommand(command)
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(command, response), nil
}

// CancelEventSubscriptions Cancel any existing event subscription.
func (socket *SocketConnection) CancelEventSubscriptions() (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	response, err := socket.sendSyncSingleLineCommand("noevents")
	if err != nil {
		return nil, err
	}
	return NewCommandResponse("noevents", response), nil
}

// AddEventFilter Add an event filter to the current set of event filters on this connection. Any of the event headers can be used as a filter.
func (socket *SocketConnection) AddEventFilter(eventHeader, valueToFilter string) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if eventHeader != "" {
		sb.WriteString("filter ")
		sb.WriteString(eventHeader)
	}
	if valueToFilter != "" {
		sb.WriteString(" ")
		sb.WriteString(valueToFilter)
	}
	response, err := socket.sendSyncSingleLineCommand(sb.String())
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(sb.String(), response), nil
}

// DeleteEventFilter Delete an event filter from the current set of event filters on this connection.
func (socket *SocketConnection) DeleteEventFilter(eventHeader, valueToFilter string) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if eventHeader != "" {
		sb.WriteString("filter delete ")
		sb.WriteString(eventHeader)
	}
	if valueToFilter != "" {
		sb.WriteString(" ")
		sb.WriteString(valueToFilter)
	}
	response, err := socket.sendSyncSingleLineCommand(sb.String())
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(sb.String(), response), nil
}

// SendEvent - Send a {@link SendEvent} command to FreeSWITCH.  This client requires that the {@link SendEvent}
//   - has a call UUID parameter.
//   - @param sendMsg a {@link SendMsg} with call UUID
//   - @return a {@link CommandResponse} with the server's response.
func (socket *SocketConnection) SendEvent(sendMsg SendEvent) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	response, err := socket.sendSyncMultiLineCommand(sendMsg.GetMsgLines())
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(sendMsg.ToString(), response), nil
}

// SendMessage - Send a {@link SendMsg} command to FreeSWITCH.  This client requires that the {@link SendMsg}
//   - has a call UUID parameter.
//   - @param sendMsg a {@link SendMsg} with call UUID
//   - @return a {@link CommandResponse} with the server's response.
func (socket *SocketConnection) SendMessage(sendMsg SendMsg) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	response, err := socket.sendSyncMultiLineCommand(sendMsg.GetMsgLines())
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(sendMsg.ToString(), response), nil
}

// SetLoggingLevel - Enable log output.
//   - @param level using the same values as in console.conf
//   - @return a {@link CommandResponse} with the server's response.
func (socket *SocketConnection) SetLoggingLevel(level string) (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if level != "" {
		sb.WriteString("log ")
		sb.WriteString(level)
	}
	response, err := socket.sendSyncSingleLineCommand(sb.String())
	if err != nil {
		return nil, err
	}
	return NewCommandResponse(sb.String(), response), nil
}

// CancelLogging - Disable any logging previously enabled with setLogLevel().
//   - @return a {@link CommandResponse} with the server's response.
func (socket *SocketConnection) CancelLogging() (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	response, err := socket.sendSyncSingleLineCommand("nolog")
	if err != nil {
		return nil, err
	}
	return NewCommandResponse("nolog", response), nil
}

// Close - Close the socket connection.
//   - @return a {@link CommandResponse} with the server's response.
func (socket *SocketConnection) Close() (*CommandResponse, error) {
	err := socket.CheckConnected()
	if err != nil {
		return nil, err
	}
	response, err := socket.sendSyncSingleLineCommand("exit")
	if err != nil {
		return nil, err
	}
	return NewCommandResponse("exit", response), nil
}

// RemoteAddr - Will return originator address known as net.RemoteAddr()
func (socket *SocketConnection) RemoteAddr() net.Addr {
	return socket.RemoteAddr()
}

func (socket *SocketConnection) CheckConnected() error {
	if socket.CanSend() {
		return nil
	}
	return errors.New("Not connected to FreeSWITCH Event Socket")
}

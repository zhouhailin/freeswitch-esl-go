package esl

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"strings"
)

func messageReceived(c *Client, m *EslMessage) error {
	contentType := m.GetContentType()
	if contentType == TEXT_EVENT_PLAIN || contentType == TEXT_EVENT_XML {
		event, err := NewEslEvent(m, true)
		if err != nil {
			return err
		}
		return handleEslEvent(c, event)
	} else {
		return handleEslMessage(contentType, c, m)
	}
}

// sendSyncSingleLineCommand - Synthesise a synchronous command/response by creating a callback object which is placed in
//   - queue and blocks waiting for another IO thread to process an incoming {@link EslMessage} and
//   - attach it to the callback.
//
// - @param channel
// - @param command single string to send
// - @return the {@link EslMessage} attached to this command's callback
func (socket *SocketConnection) sendSyncSingleLineCommand(command string) (*EslMessage, error) {
	if socket == nil {
		return nil, errors.New("connection is null.")
	}
	if isTraceEnabled() {
		logger.Tracef("sendSyncSingleLineCommand command : %s\n", command)
	}
	socket.mtx.Lock()
	defer socket.mtx.Unlock()
	_, err := socket.Writer().WriteString(command + MESSAGE_TERMINATOR)
	if err != nil {
		return nil, err
	}
	err = socket.Writer().Flush()
	if err != nil {
		return nil, err
	}
	// Block until the response is available
	return <-socket.msg, nil
}

// sendSyncMultiLineCommand - Synthesise a synchronous command/response by creating a callback object which is placed in
//   - queue and blocks waiting for another IO thread to process an incoming {@link EslMessage} and
//   - attach it to the callback.
//
// - @param channel
// - @param command List of command lines to send
// - @return the {@link EslMessage} attached to this command's callback
func (socket *SocketConnection) sendSyncMultiLineCommand(commandLines *[]string) (*EslMessage, error) {
	var sb strings.Builder
	for _, line := range *commandLines {
		sb.WriteString(line)
		sb.WriteString(LINE_TERMINATOR)
	}
	sb.WriteString(LINE_TERMINATOR)
	socket.mtx.Lock()
	defer socket.mtx.Unlock()
	_, err := socket.Writer().WriteString(sb.String())
	if err != nil {
		return nil, err
	}
	err = socket.Writer().Flush()
	if err != nil {
		return nil, err
	}
	// Block until the response is available
	return <-socket.msg, nil
}

// sendAsyncCommand - Returns the Job UUID of that the response event will have.
//   - @param channel
//   - @param command
//   - @return Job-UUID as a string
func (socket *SocketConnection) sendAsyncCommand(command string) (*string, error) {
	response, err := socket.sendSyncSingleLineCommand(command)
	if err != nil {
		return nil, err
	}
	if response.HasHeader("Job-UUID") {
		value := response.GetHeaderValue("Job-UUID")
		return &value, nil
	} else {
		return nil, errors.New("Missing Job-UUID header in bgapi response")
	}
}

func handleEslMessage(contentType string, c *Client, m *EslMessage) error {
	if isDebugEnabled() {
		logger.Debugf("Received message: %s\n", m.ToString())
	}
	switch contentType {
	case API_RESPONSE:
		if isDebugEnabled() {
			logger.Debugf("Api response received: %s\n", m.ToString())
		}
		c.msg <- m
		break
	case COMMAND_REPLY:
		if isDebugEnabled() {
			logger.Debugf("Command reply received: %s\n", m.ToString())
		}
		c.msg <- m
		break
	case AUTH_REQUEST:
		if isDebugEnabled() {
			logger.Debugf("Auth request received: %s\n", m.ToString())
		}
		go func() {
			_ = handleAuthRequest(c)
		}()
		break
	case TEXT_DISCONNECT_NOTICE:
		if isInfoEnabled() {
			logger.Infof("Disconnect notice received: %s\n", m.ToString())
		}
		return handleDisconnectionNotice(c)
	case TEXT_RUDE_REJECTION:
		if isInfoEnabled() {
			logger.Infof("Rude rejection received: %s\n", m.ToString())
		}
		return handleRudeRejection(c)
	default:
		logger.Warnf("Unexpected message content type %s", contentType)
	}
	return nil
}

func handleEslEvent(c *Client, e *EslEvent) error {
	if isDebugEnabled() {
		logger.Debugf("Received event: %s\n", e.ToString())
	}
	c.listener.eventReceived(c, e)
	return nil
}

func handleAuthRequest(c *Client) error {
	if isDebugEnabled() {
		logger.Debugf("Auth requested, sending [auth %s]\n", "*****")
	}
	response, err := c.sendSyncSingleLineCommand("auth " + c.Password)
	if err != nil {
		return err
	}
	if isDebugEnabled() {
		logger.Debugf("Auth response %s", response.ToString())
	}
	if COMMAND_REPLY == response.GetContentType() {
		c.listener.authResponseReceived(c, NewCommandResponse("auth "+c.Password, response))
		return nil
	} else {
		logger.Errorf("Bad auth response message %s\n", response)
		return errors.New("Incorrect auth response")
	}
}

func handleDisconnectionNotice(c *Client) error {
	if isDebugEnabled() {
		logger.Debug("Received disconnection notice")
	}
	c.listener.disconnected(c)
	return nil
}

func handleRudeRejection(c *Client) error {
	if isDebugEnabled() {
		logger.Debugf("Received rude rejection")
	}
	c.rudeRejection = true
	return nil
}

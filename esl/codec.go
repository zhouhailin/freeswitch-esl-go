package esl

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

// decode reads a complete ESL message from the gnet connection buffer.
// Returns io.ErrUnexpectedEOF if the buffer doesn't contain a complete message yet.
// On success, consumes the message bytes from the buffer.
func decode(c gnet.Conn, m *EslMessage) error {
	buf, err := c.Peek(-1)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Find end of headers (double LF)
	headerEndIdx := bytes.Index(buf, []byte("\n\n"))
	if headerEndIdx < 0 {
		return io.ErrUnexpectedEOF
	}

	// Parse header lines
	headerBlock := buf[:headerEndIdx]
	headerLines := bytes.Split(headerBlock, []byte("\n"))
	for _, lineBytes := range headerLines {
		headerLine := string(lineBytes)
		if isDebugEnabled() {
			logging.Debugf("read header line %s", headerLine)
		}
		if len(headerLine) == 0 {
			continue
		}
		headerParts := strings.SplitN(headerLine, ":", 2)
		if len(headerParts) < 2 {
			return errors.New("Unhandled ESL header [" + headerParts[0] + "]")
		}
		headerName := fromLiteral(headerParts[0])
		if headerName == "" {
			return errors.New("Unhandled ESL header [" + headerParts[0] + "]")
		}
		m.headers[headerName] = strings.TrimSpace(headerParts[1])
	}

	// Calculate total message length (headers + \n\n)
	totalLen := headerEndIdx + 2

	//
	// have read all headers - check for content-length
	//
	if lv := m.GetHeaderValue(CONTENT_LENGTH); lv != "" {
		if isDebugEnabled() {
			logging.Debugf("have content-length, decoding body ..")
		}
		l, err := strconv.Atoi(lv)
		if err != nil {
			logging.Errorf("Unable to get size of content-length: %s", lv)
			return err
		}
		m.contentLength = l
		if totalLen+l > len(buf) {
			return io.ErrUnexpectedEOF
		}
		bodyBytes := buf[totalLen : totalLen+l]
		if isDebugEnabled() {
			logging.Debugf("read %d body bytes", len(bodyBytes))
		}
		if len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			if bodyStr[len(bodyStr)-1] == '\n' {
				bodyStr = bodyStr[:len(bodyStr)-1]
			}
			for _, bodyLine := range strings.Split(bodyStr, LINE_TERMINATOR) {
				m.addBodyLine(bodyLine)
				if isTraceEnabled() {
					logging.Debugf("read body line %s", bodyLine)
				}
			}
		}
		totalLen += l
	}

	// Consume the full message from the buffer
	_, err = c.Next(totalLen)
	return err
}

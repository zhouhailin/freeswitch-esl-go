package esl

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/netpoll"
	"strconv"
	"strings"
)

// decode .
func decode(reader netpoll.Reader, m *EslMessage) (err error) {
	//
	// read '\n' terminated lines until reach a single '\n'
	//
	reachedDoubleLF := false
	for !reachedDoubleLF {
		// this will read or fail
		line, err := reader.Until(NEW_LINE)
		if err != nil {
			logger.Error("Decode Failure")
			return err
		}
		headerLine := string(line[:len(line)-1])
		if isDebugEnabled() {
			logger.Debugf("read header line %s\n", headerLine)
		}
		if len(headerLine) == 0 {
			reachedDoubleLF = true
		} else {
			headerParts := strings.SplitN(headerLine, ":", 2)
			headerName := fromLiteral(headerParts[0])
			if headerName == "" {
				return errors.New("Unhandled ESL header [" + headerParts[0] + "]")
			} else {
				m.headers[headerName] = strings.TrimSpace(headerParts[1])
			}
		}
	}

	//
	// have read all headers - check for content-length
	//
	if lv := m.GetHeaderValue(CONTENT_LENGTH); lv != "" {
		if isDebugEnabled() {
			logger.Debug("have content-length, decoding body ..")
		}
		l, err := strconv.Atoi(lv)
		if err != nil {
			logger.Errorf("Unable to get size of content-length: %s\n", lv)
			return err
		}
		m.contentLength = l
		if isTraceEnabled() {
			logger.Trace("Decode Body ...")
		}
		bytes, err := reader.ReadBinary(l)
		if isDebugEnabled() {
			logger.Debugf("read %d body bytes\n", len(bytes))
		}
		if err != nil {
			return err
		}
		for _, bodyLine := range strings.Split(string(bytes[:len(bytes)-1]), LINE_TERMINATOR) {
			m.addBodyLine(bodyLine)
			if isTraceEnabled() {
				logger.Tracef("read body line %s\n", bodyLine)
			}
		}
	}
	return nil
}

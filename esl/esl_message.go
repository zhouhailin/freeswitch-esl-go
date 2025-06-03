package esl

import (
	"github.com/bytedance/gopkg/util/logger"
	"strconv"
	"strings"
)

// EslMessage - Basic FreeSWITCH Event Socket messages from the server are decoded into this data object.
// - <p>
// - An ESL message is modelled as text lines.  A message always has one or more header lines, and
// - optionally may have somebody lines.
// - <p>
// - Header lines are parsed and cached in a map keyed by the {@link EslHeaders.Name} enum.  A message
// - is always expected to have a "Content-Type" header
// - <p>
// - Any Body lines are cached in a list.
type EslMessage struct {
	headers       map[Name]string
	body          []string
	contentLength int
}

// GetHeaders - All the received message headers in a map keyed by {@link EslHeaders.Name}. The string mapped value
//   - is the parsed content of the header line (ie, it does not include the header name).
//   - @return map of header values
func (m *EslMessage) GetHeaders() *map[Name]string {
	return &m.headers
}

// HasHeader - Convenience method
//   - @param headerName as a {@link EslHeaders.Name}
//   - @return true if an only if there is a header entry with the supplied header name
func (m *EslMessage) HasHeader(headerName Name) bool {
	return m.headers[headerName] != ""
}

// GetHeaderValue - Convenience method
//   - @param headerName as a {@link EslHeaders.Name}
//   - @return same as getHeaders().get( headerName )
func (m *EslMessage) GetHeaderValue(headerName Name) string {
	return m.headers[headerName]
}

// HasContentLength -Convenience method
//   - @return true if and only if a header exists with name "Content-Length"
func (m *EslMessage) HasContentLength() bool {
	return m.headers[CONTENT_LENGTH] != ""
}

// GetContentLength - Convenience method
//   - @return integer value of header with name "Content-Length"
func (m *EslMessage) GetContentLength() int {
	return m.contentLength
}

// GetContentType - Convenience method
//   - @return header value of header with name "Content-Type"
func (m *EslMessage) GetContentType() string {
	return m.headers[CONTENT_TYPE]
}

// GetBodyLines - Any received message body lines
//   - @return list with a string for each line received, may be an empty list
func (m *EslMessage) GetBodyLines() *[]string {
	return &m.body
}

// AddHeader - Used by the {@link EslMessageDecoder}.
func (m *EslMessage) addHeader(name Name, value string) {
	if isTraceEnabled() {
		logger.Tracef("adding header %s %s\n", name, value)
	}
	m.headers[name] = value
}

// AddBodyLine - Used by the {@link EslMessageDecoder}
func (m *EslMessage) addBodyLine(line string) {
	if line == "" {
		return
	}
	m.body = append(m.body, line)
}

// ToString - To String
func (m *EslMessage) ToString() string {
	var sb strings.Builder
	sb.WriteString("EslMessage: contentType=[")
	sb.WriteString(m.GetContentType())
	sb.WriteString("] headers=")
	sb.WriteString(strconv.Itoa(len(m.headers)))
	sb.WriteString("] body=")
	sb.WriteString(strconv.Itoa(len(m.body)))
	sb.WriteString(" lines.")
	return sb.String()
}

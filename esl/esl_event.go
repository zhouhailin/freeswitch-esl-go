package esl

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"net/url"
	"strconv"
	"strings"
)

// EslEvent FreeSWITCH Event Socket <strong>events</strong> are decoded into this data object.
// * <p>
// * An ESL event is modelled as a collection of text lines. An event always has several eventHeader
// * lines, and optionally may have some eventBody lines.  In addition the messageHeaders of the
// * original containing {@link EslMessage} which carried the event are also available.
// * <p>
// * The eventHeader lines are parsed and cached in a map keyed by the eventHeader name string. An event
// * is always expected to have an "Event-Name" eventHeader. Commonly used eventHeader names are coded
// * in {@link EslEventHeaderNames}
// * <p>
// * Any eventBody lines are cached in a list.
// * <p>
// * The messageHeader lines from the original message are cached in a map keyed by {@link EslHeaders.Name}.
type EslEvent struct {
	messageHeaders     *map[Name]string
	eventHeaders       map[string]string
	eventBody          []string
	decodeEventHeaders bool
}

func NewEslEvent(rawMessage *EslMessage, decodeEventHeaders bool) (*EslEvent, error) {
	event := EslEvent{
		messageHeaders: rawMessage.GetHeaders(),
		eventHeaders:   make(map[string]string, len(rawMessage.body)),
	}
	//messageHeaders: messageHeaders,
	//decodeEventHeaders: decodeEventHeaders}
	contentType := rawMessage.GetContentType()
	switch contentType {
	case TEXT_EVENT_PLAIN:
		parsePlainBody(&event, &rawMessage.body, decodeEventHeaders)
		break
	case COMMAND_REPLY:
		parsePlainBody(&event, &rawMessage.body, decodeEventHeaders)
		break
	case TEXT_EVENT_XML:
		return nil, errors.New("XML events are not yet supported")
	default:
		return nil, errors.New("Unexpected EVENT content-type: " + contentType)
	}
	return &event, nil
}

func parsePlainBody(event *EslEvent, rawBodyLines *[]string, decodeEventHeaders bool) {
	isEventBody := false
	for _, rawLine := range *rawBodyLines {
		if !isEventBody {
			headerParts := strings.SplitN(rawLine, ":", 2)
			name := strings.TrimSpace(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			if decodeEventHeaders && strings.Contains(value, "%") {
				decodedValue, err := url.QueryUnescape(value)
				if err != nil {
					logger.Warnf("Could not URL decode %s\n", value)
					event.eventHeaders[name] = value
				} else {
					if isTraceEnabled() {
						logger.Tracef("decoded from: %s\n", value)
						logger.Tracef("decoded   to: %s\n", decodedValue)
					}
					event.eventHeaders[name] = decodedValue
				}
			} else {
				if isTraceEnabled() {
					logger.Tracef("addEventHeaders : %s\n", name, value)
				}
				event.eventHeaders[name] = value
			}
			if name == "Content-Length" {
				// the remaining lines will be considered body lines
				isEventBody = true
			}
		} else {
			if len(rawLine) > 0 {
				event.eventBody = append(event.eventBody, rawLine)
			}
		}
	}
}

// GetMessageHeaders - The message headers of the original ESL message from which this event was decoded.
//   - The message headers are stored in a map keyed by {@link EslHeaders.Name}. The string mapped value
//   - is the parsed content of the header line (ie, it does not include the header name).
//   - @return map of header values
func (e *EslEvent) GetMessageHeaders() *map[Name]string {
	return e.messageHeaders
}

// GetEventHeaders - The event headers of this event. The headers are parsed and stored in a map keyed by the string
//   - name of the header, and the string mapped value is the parsed content of the event header line
//   - (ie, it does not include the header name).
//   - @return map of event header values
func (e *EslEvent) GetEventHeaders() *map[string]string {
	return &e.eventHeaders
}

// GetEventBodyLines - Any event body lines that were present in the event.
//   - @return list of decoded event body lines, may be an empty list.
func (e *EslEvent) GetEventBodyLines() *[]string {
	return &e.eventBody
}

// GetEventName - Convenience method.
//   - @return the string value of the event header "Event-Name"
func (e *EslEvent) GetEventName() string {
	return e.eventHeaders["Event-Name"]
}

// GetEventDateTimestamp - Convenience method.
//   - @return the string value of the event header "Event-Date-Timestamp"
func (e *EslEvent) GetEventDateTimestamp() string {
	return e.eventHeaders["Event-Date-Timestamp"]
}

// GetEventDateLocal - Convenience method.
//   - @return the string value of the event header "Event-Date-Local"
func (e *EslEvent) GetEventDateLocal() string {
	return e.eventHeaders["Event-Date-Local"]
}

// GetEventDateGmt - Convenience method.
//   - @return the string value of the event header "Event-Date-GMT"
func (e *EslEvent) GetEventDateGmt() string {
	return e.eventHeaders["Event-Date-GMT"]
}

// HasEventBody - Convenience method.
//   - @return true if the eventBody list is not empty
func (e *EslEvent) HasEventBody() bool {
	return len(e.eventBody) != 0
}

func (e *EslEvent) ToString() string {
	var sb strings.Builder
	sb.WriteString("EslEvent: name=[")
	sb.WriteString(e.GetEventName())
	sb.WriteString("] headers=")
	sb.WriteString(strconv.Itoa(len(*e.messageHeaders)))
	sb.WriteString(", eventHeaders=")
	sb.WriteString(strconv.Itoa(len(e.eventHeaders)))
	sb.WriteString(", eventBody=")
	sb.WriteString(strconv.Itoa(len(e.eventBody)))
	sb.WriteString(" lines.")
	return sb.String()
}

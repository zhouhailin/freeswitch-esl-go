package esl

const (
	CONTENT_TYPE   Name = "Content-Type"
	CONTENT_LENGTH Name = "Content-Length"
	REPLY_TEXT     Name = "Reply-Text"
	JOB_UUID       Name = "Job-UUID"
	SOCKET_MODE    Name = "Socket-Mode"
	CONTROL        Name = "CONTROL"
)

func fromLiteral(literal string) Name {
	switch literal {
	case "Content-Type":
		return CONTENT_TYPE
	case "Content-Length":
		return CONTENT_LENGTH
	case "Reply-Text":
		return REPLY_TEXT
	case "Job-UUID":
		return JOB_UUID
	case "Socket-Mode":
		return SOCKET_MODE
	case "CONTROL":
		return CONTROL
	default:
		return ""
	}
}

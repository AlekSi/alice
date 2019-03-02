package alice

// https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#response
type Response struct {
	Version  string          `json:"version"`
	Response ResponsePayload `json:"response"`
	Session  ResponseSession `json:"session"`
}

type ResponsePayload struct {
	Text       string           `json:"text"`
	Tts        string           `json:"tts"`
	Card       *ResponseCard    `json:"card,omitempty"`
	Buttons    []ResponseButton `json:"buttons,omitempty"`
	EndSession bool             `json:"end_session"`
}

type ResponseCard struct {
	Type   string `json:"type"`
	Header struct {
		Text string `json:"text"`
	} `json:"header,omitempty"`
	Items []struct {
		ImageID     string             `json:"image_id"`
		Title       string             `json:"title"`
		Description string             `json:"description"`
		Button      ResponseCardButton `json:"button,omitempty"`
	} `json:"items,omitempty"`
	Footer struct {
		Text   string             `json:"text"`
		Button ResponseCardButton `json:"button,omitempty"`
	} `json:"footer,omitempty"`
}

type ResponseCardButton struct {
	Text    string      `json:"text"`
	URL     string      `json:"url"`
	Payload interface{} `json:"payload,omitempty"`
}

type ResponseButton struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
	URL     string      `json:"url"`
	Hide    bool        `json:"hide"`
}

type ResponseSession struct {
	SessionID string `json:"session_id"`
	MessageID int    `json:"message_id"`
	UserID    string `json:"user_id"`
}

func NewResponse(req *Request) *Response {
	return &Response{
		Version: req.Version,
		Session: ResponseSession{
			SessionID: req.Session.SessionID,
			MessageID: req.Session.MessageID,
			UserID:    req.Session.UserID,
		},
	}
}

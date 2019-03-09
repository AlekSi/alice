package alice

type ResponseCardType string

const (
	BigImage  ResponseCardType = "BigImage"
	ItemsList ResponseCardType = "ItemsList"
)

// Response represents main response object.
// See https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#response
type Response struct {
	Response ResponsePayload `json:"response"`
	Session  ResponseSession `json:"session"`
	Version  string          `json:"version"`
}

// ResponsePayload contains response payload.
type ResponsePayload struct {
	Text       string           `json:"text"`
	Tts        string           `json:"tts,omitempty"`
	Card       *ResponseCard    `json:"card,omitempty"`
	Buttons    []ResponseButton `json:"buttons,omitempty"`
	EndSession bool             `json:"end_session"`
}

// ResponseCard contains response card.
type ResponseCard struct {
	Type ResponseCardType `json:"type"`

	// single image

	ResponseCardItem

	// multiple images

	Header *ResponseCardHeader `json:"header,omitempty"`
	Items  []ResponseCardItem  `json:"items,omitempty"`
	Footer *ResponseCardFooter `json:"footer,omitempty"`
}

type ResponseCardHeader struct {
	Text string `json:"text"`
}

type ResponseCardItem struct {
	ImageID     string              `json:"image_id,omitempty"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Button      *ResponseCardButton `json:"button,omitempty"`
}

type ResponseCardFooter struct {
	Text   string              `json:"text"`
	Button *ResponseCardButton `json:"button,omitempty"`
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

// ResponseSession contains response session.
type ResponseSession struct {
	SessionID string `json:"session_id"`
	MessageID int    `json:"message_id"`
	UserID    string `json:"user_id"`
}

// NewResponse creates new responses object and copies version and session from request.
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

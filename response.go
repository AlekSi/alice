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

	// For BigImage type.
	*ResponseCardItem `json:",omitempty"`

	// For ItemsList type.
	*ResponseCardItemsList `json:",omitempty"`
}

type ResponseCardItem struct {
	ImageID     string              `json:"image_id,omitempty"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Button      *ResponseCardButton `json:"button,omitempty"`
}

type ResponseCardItemsList struct {
	Header *ResponseCardHeader `json:"header,omitempty"`
	Items  []ResponseCardItem  `json:"items,omitempty"`
	Footer *ResponseCardFooter `json:"footer,omitempty"`
}

type ResponseCardHeader struct {
	Text string `json:"text"`
}

type ResponseCardFooter struct {
	Text   string              `json:"text"`
	Button *ResponseCardButton `json:"button,omitempty"`
}

type ResponseCardButton struct {
	Text    string      `json:"text"`
	URL     string      `json:"url,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type ResponseButton struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
	URL     string      `json:"url,omitempty"`
	Hide    bool        `json:"hide,omitempty"`
}

// ResponseSession contains response session.
type ResponseSession struct {
	SessionID string `json:"session_id"`
	MessageID int    `json:"message_id"`
	UserID    string `json:"user_id"`
}

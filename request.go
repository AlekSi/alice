package alice

type RequestPayloadType string

const (
	SimpleUtterance RequestPayloadType = "SimpleUtterance"
	ButtonPressed   RequestPayloadType = "ButtonPressed"
)

// https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#request
type Request struct {
	Version string         `json:"version"`
	Meta    RequestMeta    `json:"meta"`
	Request RequestPayload `json:"request"`
	Session RequestSession `json:"session"`
}

type RequestMeta struct {
	Locale     string                 `json:"locale"`
	Timezone   string                 `json:"timezone"`
	ClientID   string                 `json:"client_id"`
	Interfaces map[string]interface{} `json:"interfaces"`
}

type RequestPayload struct {
	Command           string             `json:"command"`
	OriginalUtterance string             `json:"original_utterance"`
	Type              RequestPayloadType `json:"type"`
	Markup            struct {
		DangerousContext bool `json:"dangerous_context"`
	} `json:"markup"`
	Payload interface{} `json:"payload,omitempty"`
	NLU     RequestNLU  `json:"nlu"`
}

type RequestNLU struct {
	Tokens   []string      `json:"tokens"`
	Entities []interface{} `json:"entities"` // TODO
}

type RequestSession struct {
	New       bool   `json:"new"`
	MessageID int    `json:"message_id"`
	SessionID string `json:"session_id"`
	SkillID   string `json:"skill_id"`
	UserID    string `json:"user_id"`
}

func (m RequestMeta) HasScreen() bool {
	return m.Interfaces["screen"] != nil
}

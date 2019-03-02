package alice

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseDecode(t *testing.T) {
	// from https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#response

	t.Run("NoImage", func(t *testing.T) {
		// from https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#request
		b := []byte(`
		{
			"response": {
				"text": "Здравствуйте! Это мы, хороводоведы.",
				"tts": "Здравствуйте! Это мы, хоров+одо в+еды.",
				"buttons": [
					{
						"title": "Надпись на кнопке",
						"payload": {},
						"url": "https://example.com/",
						"hide": true
					}
				],
				"end_session": false
			},
			"session": {
				"session_id": "2eac4854-fce721f3-b845abba-20d60",
				"message_id": 4,
				"user_id": "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"
			},
			"version": "1.0"
		}`)
		d := json.NewDecoder(bytes.NewReader(b))
		d.DisallowUnknownFields()
		var resp Response
		if err := d.Decode(&resp); err != nil {
			t.Fatal(err)
		}
		expected := []ResponseButton{{
			Title:   "Надпись на кнопке",
			Payload: map[string]interface{}{},
			URL:     "https://example.com/",
			Hide:    true,
		}}
		assert.Equal(t, expected, resp.Response.Buttons)
	})
}

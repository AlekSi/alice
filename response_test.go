package alice

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	// from https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#response

	t.Run("NoImage", func(t *testing.T) {
		b := []byte(strings.TrimSpace(`
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
		}`))
		var actual Response

		t.Run("Decode", func(t *testing.T) {
			d := json.NewDecoder(bytes.NewReader(b))
			d.DisallowUnknownFields()
			err := d.Decode(&actual)
			require.NoError(t, err, "%#v", err)
			expected := Response{
				Version: "1.0",
				Response: ResponsePayload{
					Text: "Здравствуйте! Это мы, хороводоведы.",
					Tts:  "Здравствуйте! Это мы, хоров+одо в+еды.",
					Buttons: []ResponseButton{{
						Title:   "Надпись на кнопке",
						Payload: map[string]interface{}{},
						URL:     "https://example.com/",
						Hide:    true,
					}},
				},
				Session: ResponseSession{
					SessionID: "2eac4854-fce721f3-b845abba-20d60",
					MessageID: 4,
					UserID:    "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC",
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("Encode", func(t *testing.T) {
			actualB, err := json.MarshalIndent(actual, "\t\t", "\t")
			require.NoError(t, err)
			assert.Equal(t, strings.Split(string(b), "\n"), strings.Split(string(actualB), "\n"))
		})
	})

	t.Run("BigImage", func(t *testing.T) {
		b := []byte(strings.TrimSpace(`
		{
			"response": {
				"text": "Здравствуйте! Это мы, хороводоведы.",
				"tts": "Здравствуйте! Это мы, хоров+одо в+еды.",
				"card": {
					"type": "BigImage",
					"image_id": "1027858/46r960da47f60207e924",
					"title": "Заголовок для изображения",
					"description": "Описание изображения.",
					"button": {
						"text": "Надпись на кнопке",
						"url": "http://example.com/",
						"payload": {}
					}
				},
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
		}`))
		var actual Response

		t.Run("Decode", func(t *testing.T) {
			d := json.NewDecoder(bytes.NewReader(b))
			d.DisallowUnknownFields()
			err := d.Decode(&actual)
			require.NoError(t, err, "%#v", err)
			expected := Response{
				Version: "1.0",
				Response: ResponsePayload{
					Text: "Здравствуйте! Это мы, хороводоведы.",
					Tts:  "Здравствуйте! Это мы, хоров+одо в+еды.",
					Card: &ResponseCard{
						Type: BigImage,
						ResponseCardItem: &ResponseCardItem{
							ImageID:     "1027858/46r960da47f60207e924",
							Title:       "Заголовок для изображения",
							Description: "Описание изображения.",
							Button: &ResponseCardButton{
								Text:    "Надпись на кнопке",
								URL:     "http://example.com/",
								Payload: map[string]interface{}{},
							},
						},
					},
					Buttons: []ResponseButton{{
						Title:   "Надпись на кнопке",
						Payload: map[string]interface{}{},
						URL:     "https://example.com/",
						Hide:    true,
					}},
				},
				Session: ResponseSession{
					SessionID: "2eac4854-fce721f3-b845abba-20d60",
					MessageID: 4,
					UserID:    "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC",
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("Encode", func(t *testing.T) {
			actualB, err := json.MarshalIndent(actual, "\t\t", "\t")
			require.NoError(t, err)
			assert.Equal(t, strings.Split(string(b), "\n"), strings.Split(string(actualB), "\n"))
		})
	})

	t.Run("ItemsList", func(t *testing.T) {
		b := []byte(strings.TrimSpace(`
		{
			"response": {
				"text": "Здравствуйте! Это мы, хороводоведы.",
				"tts": "Здравствуйте! Это мы, хоров+одо в+еды.",
				"card": {
					"type": "ItemsList",
					"header": {
						"text": "Заголовок галереи изображений"
					},
					"items": [
						{
							"image_id": "image_id",
							"title": "Заголовок для изображения.",
							"description": "Описание изображения.",
							"button": {
								"text": "Надпись на кнопке",
								"url": "http://example.com/",
								"payload": {}
							}
						}
					],
					"footer": {
						"text": "Текст блока под изображением.",
						"button": {
							"text": "Надпись на кнопке",
							"url": "https://example.com/",
							"payload": {}
						}
					}
				},
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
		}`))
		var actual Response

		t.Run("Decode", func(t *testing.T) {
			d := json.NewDecoder(bytes.NewReader(b))
			d.DisallowUnknownFields()
			err := d.Decode(&actual)
			require.NoError(t, err, "%#v", err)
			expected := Response{
				Version: "1.0",
				Response: ResponsePayload{
					Text: "Здравствуйте! Это мы, хороводоведы.",
					Tts:  "Здравствуйте! Это мы, хоров+одо в+еды.",
					Card: &ResponseCard{
						Type: ItemsList,
						ResponseCardItemsList: &ResponseCardItemsList{
							Header: &ResponseCardHeader{
								Text: "Заголовок галереи изображений",
							},
							Items: []ResponseCardItem{
								{
									ImageID:     "image_id",
									Title:       "Заголовок для изображения.",
									Description: "Описание изображения.",
									Button: &ResponseCardButton{
										Text:    "Надпись на кнопке",
										URL:     "http://example.com/",
										Payload: map[string]interface{}{},
									},
								},
							},
							Footer: &ResponseCardFooter{
								Text: "Текст блока под изображением.",
								Button: &ResponseCardButton{
									Text:    "Надпись на кнопке",
									URL:     "https://example.com/",
									Payload: map[string]interface{}{},
								},
							},
						},
					},
					Buttons: []ResponseButton{{
						Title:   "Надпись на кнопке",
						Payload: map[string]interface{}{},
						URL:     "https://example.com/",
						Hide:    true,
					}},
				},
				Session: ResponseSession{
					SessionID: "2eac4854-fce721f3-b845abba-20d60",
					MessageID: 4,
					UserID:    "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC",
				},
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("Encode", func(t *testing.T) {
			actualB, err := json.MarshalIndent(actual, "\t\t", "\t")
			require.NoError(t, err)
			assert.Equal(t, strings.Split(string(b), "\n"), strings.Split(string(actualB), "\n"))
		})
	})
}

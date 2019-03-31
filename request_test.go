package alice

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestDecode(t *testing.T) {
	// from https://tech.yandex.ru/dialogs/alice/doc/protocol-docpage/#request
	b := []byte(`
	{
		"meta": {
			"locale": "ru-RU",
			"timezone": "Europe/Moscow",
			"client_id": "ru.yandex.searchplugin/5.80 (Samsung Galaxy; Android 4.4)",
			"interfaces": {
				"screen": {}
			}
		},
		"request": {
			"command": "закажи пиццу на улицу льва толстого, 16 на завтра",
			"original_utterance": "закажи пиццу на улицу льва толстого, 16 на завтра",
			"type": "SimpleUtterance",
			"markup": {
				"dangerous_context": true
			},
			"payload": {},
			"nlu": {
				"tokens": [
					"закажи",
					"пиццу",
					"на",
					"льва",
					"толстого",
					"16",
					"на",
					"завтра"
				],
				"entities": [
					{
						"tokens": {
							"start": 2,
							"end": 6
						},
						"type": "YANDEX.GEO",
						"value": {
							"house_number": "16",
							"street": "льва толстого"
						}
					},
					{
						"tokens": {
							"start": 3,
							"end": 5
						},
						"type": "YANDEX.FIO",
						"value": {
							"first_name": "лев",
							"last_name": "толстой"
						}
					},
					{
						"tokens": {
							"start": 5,
							"end": 6
						},
						"type": "YANDEX.NUMBER",
						"value": 16
					},
					{
						"tokens": {
							"start": 6,
							"end": 8
						},
						"type": "YANDEX.DATETIME",
						"value": {
							"day": 1,
							"day_is_relative": true
						}
					}
				]
			}
		},
		"session": {
			"new": true,
			"message_id": 4,
			"session_id": "2eac4854-fce721f3-b845abba-20d60",
			"skill_id": "3ad36498-f5rd-4079-a14b-788652932056",
			"user_id": "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"
		},
		"version": "1.0"
	}`)

	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	var req Request
	require.NoError(t, d.Decode(&req))
	assert.True(t, req.Meta.HasScreen())
	assert.Equal(t, []string{"закажи", "пиццу", "на", "льва", "толстого", "16", "на", "завтра"}, req.Request.NLU.Tokens)
	require.Len(t, req.Request.NLU.Entities, 4)

	geo := req.Request.NLU.Entities[0]
	assert.Equal(t, 2, geo.Tokens.Start)
	assert.Equal(t, 6, geo.Tokens.End)
	assert.Equal(t, EntityYandexGeo, geo.Type)
	assert.Equal(t, &YandexGeo{Street: "льва толстого", HouseNumber: "16"}, geo.YandexGeo())

	fio := req.Request.NLU.Entities[1]
	assert.Equal(t, 3, fio.Tokens.Start)
	assert.Equal(t, 5, fio.Tokens.End)
	assert.Equal(t, EntityYandexFio, fio.Type)
	assert.Equal(t, &YandexFio{FirstName: "лев", LastName: "толстой"}, fio.YandexFio())

	num := req.Request.NLU.Entities[2]
	assert.Equal(t, 5, num.Tokens.Start)
	assert.Equal(t, 6, num.Tokens.End)
	assert.Equal(t, EntityYandexNumber, num.Type)
	assert.Equal(t, "16", num.YandexNumber().String())

	dt := req.Request.NLU.Entities[3]
	assert.Equal(t, 6, dt.Tokens.Start)
	assert.Equal(t, 8, dt.Tokens.End)
	assert.Equal(t, EntityYandexDateTime, dt.Type)
	assert.Equal(t, &YandexDateTime{Day: 1, DayIsRelative: true}, dt.YandexDateTime())
}

package alice

import (
	"encoding/json"
	"strconv"
)

// EntityType represents entity type.
type EntityType string

// Entity types.
const (
	EntityYandexFio      EntityType = "YANDEX.FIO"
	EntityYandexGeo      EntityType = "YANDEX.GEO"
	EntityYandexDateTime EntityType = "YANDEX.DATETIME"
	EntityYandexNumber   EntityType = "YANDEX.NUMBER"
)

// Entity represents NLU-extracted named entity.
// See https://yandex.ru/dev/dialogs/alice/doc/nlu-docpage/.
type Entity struct {
	Tokens EntityTokens `json:"tokens"`
	Type   EntityType   `json:"type"`
	Value  interface{}  `json:"value"`
}

// EntityTokens represents the place of extracted named entity in tokens slice.
type EntityTokens struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// YandexFio represents extracted full name.
type YandexFio struct {
	FirstName      string
	PatronymicName string
	LastName       string
}

// YandexFio extracts full name from entity, or nil.
func (e *Entity) YandexFio() *YandexFio {
	if e.Type != EntityYandexFio {
		return nil
	}
	v, _ := e.Value.(map[string]interface{})
	if v == nil {
		return nil
	}

	f, _ := v["first_name"].(string)
	p, _ := v["patronymic_name"].(string)
	l, _ := v["last_name"].(string)

	return &YandexFio{
		FirstName:      f,
		PatronymicName: p,
		LastName:       l,
	}
}

// YandexGeo represents extracted address.
type YandexGeo struct {
	Country     string
	City        string
	Street      string
	HouseNumber string
	Airport     string
}

// YandexGeo extracts address from entity, or nil.
func (e *Entity) YandexGeo() *YandexGeo {
	if e.Type != EntityYandexGeo {
		return nil
	}
	v, _ := e.Value.(map[string]interface{})
	if v == nil {
		return nil
	}

	co, _ := v["country"].(string)
	ci, _ := v["city"].(string)
	s, _ := v["street"].(string)
	h, _ := v["house_number"].(string)
	a, _ := v["airport"].(string)

	return &YandexGeo{
		Country:     co,
		City:        ci,
		Street:      s,
		HouseNumber: h,
		Airport:     a,
	}
}

// YandexDateTime represents extracted date and/or time.
type YandexDateTime struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int

	YearIsRelative   bool
	MonthIsRelative  bool
	DayIsRelative    bool
	HourIsRelative   bool
	MinuteIsRelative bool
}

// YandexDateTime extracts date and/or time from entity, or nil.
func (e *Entity) YandexDateTime() *YandexDateTime {
	if e.Type != EntityYandexDateTime {
		return nil
	}
	v, _ := e.Value.(map[string]interface{})
	if v == nil {
		return nil
	}

	keys := []string{"year", "month", "day", "hour", "minute"}

	// extract absolute values
	abs := make(map[string]int)
	for _, k := range keys {
		switch v := v[k].(type) {
		case json.Number:
			i64, _ := v.Int64()
			abs[k] = int(i64)
		case float64:
			abs[k] = int(v)
		}
	}

	// extract relative flags
	rel := make(map[string]bool)
	for _, k := range keys {
		rel[k], _ = v[k+"_is_relative"].(bool)
	}

	return &YandexDateTime{
		Year:   abs["year"],
		Month:  abs["month"],
		Day:    abs["day"],
		Hour:   abs["hour"],
		Minute: abs["minute"],

		YearIsRelative:   rel["year"],
		MonthIsRelative:  rel["month"],
		DayIsRelative:    rel["day"],
		HourIsRelative:   rel["hour"],
		MinuteIsRelative: rel["minute"],
	}
}

// YandexNumber represents extracted integer or float number.
type YandexNumber struct {
	json.Number
}

// YandexNumber extracts integer or float number from entity, or nil.
func (e *Entity) YandexNumber() *YandexNumber {
	if e.Type != EntityYandexNumber {
		return nil
	}

	switch v := e.Value.(type) {
	case json.Number:
		return &YandexNumber{v}
	case float64:
		return &YandexNumber{json.Number(strconv.FormatFloat(v, 'f', -1, 64))}
	default:
		return nil
	}
}

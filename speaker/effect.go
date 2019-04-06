// Package speaker provides various sounds and effects.
package speaker

import (
	"fmt"
)

// Effect defines an effect you can apply to text.
// See https://yandex.ru/dev/dialogs/alice/doc/speech-effects-docpage/.
type Effect string

// Those contants define various effects.
const (
	BehindTheWall = Effect("behind_the_wall")
	Hamster       = Effect("hamster")
	Megaphone     = Effect("megaphone")
	PitchDown     = Effect("pitch_down")
	Psychodelic   = Effect("psychodelic")
	Pulse         = Effect("pulse")
	TrainAnnounce = Effect("train_announce")
)

// Apply wraps given text with effect.
func (e Effect) Apply(text string) string {
	return fmt.Sprintf(`<speaker effect="%s">%s<speaker effect="-">`, e, text)
}

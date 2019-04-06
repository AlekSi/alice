package speaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEffect(t *testing.T) {
	expected := `<speaker effect="behind_the_wall">Behind!<speaker effect="-"><speaker audio="alice-sounds-game-win-1.opus">`
	assert.Equal(t, expected, BehindTheWall.Apply("Behind!")+Chainsaw1)
}

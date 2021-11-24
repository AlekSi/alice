package resources

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("-short is passed, skipping integration test")
	}

	skillID := os.Getenv("ALICE_TEST_SKILL_ID")
	oAuthToken := os.Getenv("ALICE_TEST_OAUTH_TOKEN")
	if skillID == "" || oAuthToken == "" {
		t.Skip("`ALICE_TEST_SKILL_ID` or `ALICE_TEST_OAUTH_TOKEN` is not set, skipping integration test")
	}

	c := Client{
		SkillID:       skillID,
		OAuthToken:    oAuthToken,
		Indent:        true,
		StrictDecoder: true,
	}

	t.Run("Status", func(t *testing.T) {
		c.Debugf = t.Logf

		status, err := c.Status()
		require.NoError(t, err)
		assert.Equal(t, 104857600, status.Images.Quota.Total)
		assert.Equal(t, 1073741824, status.Sounds.Quota.Total)
	})

	t.Run("Sound", func(t *testing.T) {
		t.Run("UploadSoundFile", func(t *testing.T) {
			c.Debugf = t.Logf

			sound, err := c.UploadSoundFile(filepath.Join("..", "testdata", "go.wav"))
			require.NoError(t, err)
			require.NotEmpty(t, sound)
			assert.NotEmpty(t, skillID, sound.ID)
			assert.Equal(t, skillID, sound.SkillID)
			assert.Empty(t, sound.Size)
			assert.Equal(t, "go.wav", sound.OriginalName)
			assert.WithinDuration(t, time.Now(), sound.CreatedAt, 5*time.Second)
			assert.False(t, sound.IsProcessed)
			assert.Nil(t, sound.Error)

			sounds, err := c.ListSounds()
			require.NoError(t, err)
			assert.Empty(t, sounds)
		})
	})
}

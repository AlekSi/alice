package resources

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("-short is passed, skipping integration test")
	}

	skill := os.Getenv("ALICE_TEST_SKILL_ID")
	token := os.Getenv("ALICE_TEST_OAUTH_TOKEN")
	if skill == "" || token == "" {
		t.Skip("`ALICE_TEST_SKILL_ID` or `ALICE_TEST_OAUTH_TOKEN` is not set, skipping integration test")
	}

	c := Client{
		OAuthToken: token,
		SkillID:    skill,
	}

	t.Run("Status", func(t *testing.T) {
		status, err := c.Status()
		require.NoError(t, err)
		assert.Equal(t, 104857600, status.Total)
	})
}

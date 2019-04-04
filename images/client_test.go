package images

import (
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("-short is passed, skipping integration test")
	}

	token := os.Getenv("ALICE_TEST_OAUTH_TOKEN")
	skill := os.Getenv("ALICE_TEST_SKILL_ID")
	if token == "" || skill == "" {
		t.Skip("`ALICE_TEST_OAUTH_TOKEN` or `ALICE_TEST_SKILL_ID` is not set, skipping integration test")
	}

	c := Client{
		OAuthToken: token,
		SkillID:    skill,
	}

	t.Run("Status", func(t *testing.T) {
		status, err := c.Status()
		if err != nil {
			t.Fatal(err)
		}
		if status.Total != 104857600 {
			t.Errorf("unexpected total %d", status.Total)
		}
	})
}

package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	SkillID    string
	OAuthToken string
	HTTPClient *http.Client
}

func (c *Client) http() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

type StatusResponse struct {
	Total int
	Used  int
}

func (c *Client) Status() (*StatusResponse, error) {
	req, err := http.NewRequest("GET", "https://dialogs.yandex.net/api/v1/status", nil)
	if err != nil {
		return nil, err
	}
	if c.OAuthToken != "" {
		req.Header.Set("Authorization", "OAuth "+c.OAuthToken)
	}

	resp, err := c.http().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var res struct {
		Images struct {
			Quota StatusResponse
		}
	}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Images.Quota, nil
}

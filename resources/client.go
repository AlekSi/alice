package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/AlekSi/alice"
)

type Quota struct {
	Total int
	Used  int
}

type Sound struct {
	ID           string
	SkillID      string
	Size         *int
	OriginalName string
	CreatedAt    time.Time
	IsProcessed  bool
	Error        *string
}

type Client struct {
	SkillID    string
	OAuthToken string
	HTTPClient *http.Client

	// debugging options
	Debugf        alice.Printf // debug logger
	Indent        bool         // indent requests and responses
	StrictDecoder bool         // disallow unexpected fields in responses
}

func (c *Client) do(req *http.Request, respBody interface{}) error {
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	var jsonRequst bool
	if c.OAuthToken != "" {
		req.Header.Set("Authorization", "OAuth "+c.OAuthToken)
	}
	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		jsonRequst = true
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if c.Debugf != nil {
		if c.Indent && jsonRequst {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}

			var body bytes.Buffer
			if err = json.Indent(&body, b, "", "  "); err != nil {
				return err
			}
			req.Body = ioutil.NopCloser(&body)
			req.ContentLength = int64(body.Len())
			req.TransferEncoding = nil
		}

		b, err := httputil.DumpRequestOut(req, jsonRequst)
		if err != nil {
			return err
		}
		c.debugf("Request:\n%s", b)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if c.Debugf != nil {
		if c.Indent {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			var body bytes.Buffer
			if err = json.Indent(&body, b, "", "  "); err != nil {
				return err
			}
			resp.Body = ioutil.NopCloser(&body)
			resp.ContentLength = int64(body.Len())
			resp.TransferEncoding = nil
		}

		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		c.debugf("Response:\n%s", b)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)
	if c.StrictDecoder {
		d.DisallowUnknownFields()
	}
	return d.Decode(&respBody)
}

func (c *Client) debugf(format string, a ...interface{}) {
	if c.Debugf != nil {
		c.Debugf(format, a...)
	}
}

type StatusResponse struct {
	Images struct {
		Quota Quota
	}
	Sounds struct {
		Quota Quota
	}
}

func (c *Client) Status() (*StatusResponse, error) {
	req, err := http.NewRequest("GET", "https://dialogs.yandex.net/api/v1/status", nil)
	if err != nil {
		return nil, err
	}

	var res StatusResponse
	if err = c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UploadSound(name string, r io.Reader) (*Sound, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, r); err != nil {
		return nil, err
	}
	if err = mw.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://dialogs.yandex.net/api/v1/skills/"+c.SkillID+"/sounds", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", mw.FormDataContentType())

	var res struct {
		Sound Sound
	}
	if err = c.do(req, &res); err != nil {
		return nil, err
	}
	return &res.Sound, nil
}

func (c *Client) UploadSoundFile(filename string) (*Sound, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck

	return c.UploadSound(filepath.Base(filename), f)
}

func (c *Client) ListSounds() ([]Sound, error) {
	req, err := http.NewRequest("GET", "https://dialogs.yandex.net/api/v1/skills/"+c.SkillID+"/sounds", nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Sounds []Sound
		Total  int
	}
	if err = c.do(req, &res); err != nil {
		return nil, err
	}
	return res.Sounds, nil
}

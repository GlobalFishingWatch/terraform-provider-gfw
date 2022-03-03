package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GFWClient struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient -
func NewClient(host, token string) (*GFWClient, error) {
	c := GFWClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		HostURL: host,
	}

	c.Token = token
	return &c, nil
}

func (c *GFWClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("permissions", `[{ "action": "read-all", "type": "entity", "value": "action" },{ "action": "create-all", "type": "entity", "value": "action" }]`)
	req.Header.Set("user", `{"id":-1}`)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

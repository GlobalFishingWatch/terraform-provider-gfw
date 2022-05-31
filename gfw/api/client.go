package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

func valuesToRaw(query url.Values) string {
	rawQuery := ""
	count := 0
	for k, v := range query {
		if count == 0 {
			rawQuery = fmt.Sprintf("%s=%s", k, url.QueryEscape(strings.Join(v, ",")))
		} else {
			rawQuery = fmt.Sprintf("%s&%s=%s", rawQuery, k, url.QueryEscape(strings.Join(v, ",")))
		}
		count++
	}
	return rawQuery
}

func (c *GFWClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	query := req.URL.Query()
	query.Set("cache", "false")

	req.URL.RawQuery = valuesToRaw(query)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		appError := AppError{}
		err = json.Unmarshal(body, &appError)
		if err != nil {
			return nil, err
		}
		return nil, appError
	}

	return body, err
}

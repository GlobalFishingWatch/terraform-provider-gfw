package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const DATAVIEW_PATH = "dataviews"

func (c *GFWClient) GetDataviews() (*[]Dataview, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, DATAVIEW_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dataviews := Pagination[Dataview]{}
	err = json.Unmarshal(body, &dataviews)
	if err != nil {
		return nil, err
	}

	return &dataviews.Entries, nil
}

func (c *GFWClient) GetDataview(id string) (*Dataview, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, DATAVIEW_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dataview := Dataview{}
	err = json.Unmarshal(body, &dataview)
	if err != nil {
		return nil, err
	}

	return &dataview, nil
}

func (c *GFWClient) DeleteDataview(id string) (*Dataview, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, DATAVIEW_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	dataview := Dataview{}
	err = json.Unmarshal(body, &dataview)
	if err != nil {
		return nil, err
	}

	return &dataview, nil
}

func (c *GFWClient) UpdateDataview(id string, dataview CreateDataview) error {

	bodyReq, err := json.Marshal(dataview)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s/%s", c.HostURL, DATAVIEW_PATH, id), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err

}

func (c *GFWClient) CreateDataview(dataview CreateDataview) (*Dataview, error) {
	exists, err := c.checkExistDataview(dataview.Slug)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(dataview)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, DATAVIEW_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDataview := Dataview{}
	err = json.Unmarshal(body, &newDataview)
	if err != nil {
		return nil, err
	}

	return &newDataview, nil
}

func (c *GFWClient) checkExistDataview(id string) (*Dataview, error) {
	exists, err := c.GetDataview(id)
	if err != nil {
		if re, ok := err.(AppError); ok {
			if re.Code == NotFoundCode {
				return nil, nil
			}
		}
		return nil, err
	}

	return exists, nil
}

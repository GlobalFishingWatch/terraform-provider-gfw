package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const ACTION_PATH = "actions"

func (c *GFWClient) GetActions() (*[]Action, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, ACTION_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	actions := []Action{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return nil, err
	}

	return &actions, nil
}

func (c *GFWClient) GetAction(id string) (*Action, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, ACTION_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	action := Action{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) DeleteAction(id string) (*Action, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, ACTION_PATH, id), nil)
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
	action := Action{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) CreateAction(action CreateAction) (*Action, error) {
	exists, err := c.checkExistAction(action.Name)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, ACTION_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newAction := Action{}
	err = json.Unmarshal(body, &newAction)
	if err != nil {
		return nil, err
	}

	return &newAction, nil
}

func (c *GFWClient) checkExistAction(name string) (*Action, error) {
	actions, err := c.GetActions()
	if err != nil {
		return nil, err
	}
	for _, a := range *actions {
		if a.Name == name {
			return &a, nil
		}
	}
	return nil, nil
}

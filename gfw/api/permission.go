package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const PERMISSION_PATH = "permissions"

func (c *GFWClient) GetPermissions() (*[]Permission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, PERMISSION_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	permissions := []Permission{}
	err = json.Unmarshal(body, &permissions)
	if err != nil {
		return nil, err
	}

	return &permissions, nil
}

func (c *GFWClient) GetPermission(id string) (*Permission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, PERMISSION_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	permission := Permission{}
	err = json.Unmarshal(body, &permission)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (c *GFWClient) DeletePermission(id string) (*Permission, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, PERMISSION_PATH, id), nil)
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
	permission := Permission{}
	err = json.Unmarshal(body, &permission)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (c *GFWClient) CreatePermission(permission CreatePermission) (*Permission, error) {
	exists, err := c.checkExistPermission(permission.Resource, permission.Action)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(permission)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, PERMISSION_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newPermission := Permission{}
	err = json.Unmarshal(body, &newPermission)
	if err != nil {
		return nil, err
	}

	return &newPermission, nil
}

func (c *GFWClient) checkExistPermission(resourceID, actionID int) (*Permission, error) {
	permissions, err := c.GetPermissions()
	if err != nil {
		return nil, err
	}
	for _, p := range *permissions {
		if p.Resource.ID == resourceID && p.Action.ID == actionID {
			return &p, nil
		}
	}
	return nil, nil
}

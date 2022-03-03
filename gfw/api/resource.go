package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const RESOURCE_PATH = "resources"

func (c *GFWClient) GetResources() (*[]Resource, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, RESOURCE_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resources := []Resource{}
	err = json.Unmarshal(body, &resources)
	if err != nil {
		return nil, err
	}

	return &resources, nil
}

func (c *GFWClient) GetResource(id string) (*Resource, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, RESOURCE_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resource := Resource{}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (c *GFWClient) DeleteResource(id string) (*Resource, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, RESOURCE_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resource := Resource{}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (c *GFWClient) CreateResource(resource CreateResource) (*Resource, error) {
	exists, err := c.checkExistResource(resource.Type, resource.Value)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, RESOURCE_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newResource := Resource{}
	err = json.Unmarshal(body, &newResource)
	if err != nil {
		return nil, err
	}

	return &newResource, nil
}

func (c *GFWClient) checkExistResource(rType, rValue string) (*Resource, error) {
	resources, err := c.GetResources()
	if err != nil {
		return nil, err
	}
	for _, a := range *resources {
		if a.Type == rType && a.Value == rValue {
			return &a, nil
		}
	}
	return nil, nil
}

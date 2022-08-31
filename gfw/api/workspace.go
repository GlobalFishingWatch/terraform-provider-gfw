package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/iancoleman/strcase"
)

const WORKSPACE_PATH = "workspaces"

func (c *GFWClient) GetWorkspaces() (*[]Workspace, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, WORKSPACE_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspaces := Pagination[Workspace]{}
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return &workspaces.Entries, nil
}

func (c *GFWClient) GetWorkspace(id string) (*Workspace, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, WORKSPACE_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}
	workspace.Public = strings.HasSuffix(workspace.ID, "-public")
	return &workspace, nil
}

func (c *GFWClient) DeleteWorkspace(id string) (*Workspace, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, WORKSPACE_PATH, id), nil)
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
	workspace := Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *GFWClient) UpdateWorkspace(id string, workspace CreateWorkspace) error {

	bodyReq, err := json.Marshal(workspace)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s/%s", c.HostURL, WORKSPACE_PATH, id), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err

}

func (c *GFWClient) CreateWorkspace(workspace CreateWorkspace) (*Workspace, error) {
	id := workspace.ID
	if id == "" {
		id = strcase.ToSnake(workspace.Name)
		if workspace.Public {
			id = fmt.Sprintf("%s-public", id)
		}
	}
	exists, err := c.checkExistWorkspace(id)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, WORKSPACE_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newWorkspace := Workspace{}
	err = json.Unmarshal(body, &newWorkspace)
	if err != nil {
		return nil, err
	}

	return &newWorkspace, nil
}

func (c *GFWClient) checkExistWorkspace(id string) (*Workspace, error) {
	exists, err := c.GetWorkspace(id)
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

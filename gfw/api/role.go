package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
)

const ROLE_PATH = "roles"

func (c *GFWClient) GetRoles() (*[]Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, ROLE_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	actions := []Role{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return nil, err
	}

	return &actions, nil
}

func (c *GFWClient) GetRole(id string) (*Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, ROLE_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	action := Role{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) DeleteRole(id string) (*Role, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, ROLE_PATH, id), nil)
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
	action := Role{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) CreateRole(action CreateRole) (*Role, error) {
	exists, err := c.checkExistRole(action.Name)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, ROLE_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newRole := Role{}
	err = json.Unmarshal(body, &newRole)
	if err != nil {
		return nil, err
	}

	return &newRole, nil
}

func (c *GFWClient) checkExistRole(name string) (*Role, error) {
	actions, err := c.GetRoles()
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

func (c *GFWClient) CreateRolePermissions(rolePerm CreateRolePermissions) error {
	// obtain Role
	role, err := c.GetRole(strconv.Itoa(rolePerm.RoleID))
	if err != nil {
		return err
	}
	var permsToCreate []int
	var permsToDelete []int
	for _, p := range role.Permissions {
		if !utils.Exists(p.ID, rolePerm.Permissions) {
			permsToDelete = append(permsToDelete, p.ID)
		}
	}

	for _, p := range rolePerm.Permissions {
		if !existsInPermission(p, role.Permissions) {
			permsToCreate = append(permsToCreate, p)
		}
	}

	for _, createPermId := range permsToCreate {
		_, err := c.AddPermissionInRole(rolePerm.RoleID, createPermId)
		if err != nil {
			return err
		}
	}
	for _, delPermId := range permsToDelete {
		_, err := c.DeletePermissionInRole(rolePerm.RoleID, delPermId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *GFWClient) DeleteRolePermissions(roleId string) error {
	role, err := c.GetRole(roleId)
	if err != nil {
		return err
	}
	for _, p := range role.Permissions {
		_, err := c.DeletePermissionInRole(role.ID, p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GFWClient) AddPermissionInRole(roleId, permissionId int) (*Role, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%d/permission/%d", c.HostURL, ROLE_PATH, roleId, permissionId), nil)
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newRole := Role{}
	err = json.Unmarshal(body, &newRole)
	if err != nil {
		return nil, err
	}

	return &newRole, nil
}

func (c *GFWClient) DeletePermissionInRole(roleId, permissionId int) (*Role, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%d/permission/%d", c.HostURL, ROLE_PATH, roleId, permissionId), nil)
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newRole := Role{}
	err = json.Unmarshal(body, &newRole)
	if err != nil {
		return nil, err
	}

	return &newRole, nil
}

func existsInPermission(i int, array []Permission) bool {

	for _, v := range array {
		if i == v.ID {
			return true
		}
	}
	return false
}

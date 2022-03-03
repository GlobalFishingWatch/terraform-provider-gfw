package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
)

const USER_GROUP_PATH = "user-groups"

func (c *GFWClient) GetUserGroups() (*[]UserGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, USER_GROUP_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	actions := []UserGroup{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return nil, err
	}

	return &actions, nil
}

func (c *GFWClient) GetUserGroup(id string) (*UserGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, USER_GROUP_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	action := UserGroup{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) DeleteUserGroup(id string) (*UserGroup, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, USER_GROUP_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	action := UserGroup{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (c *GFWClient) CreateUserGroup(action CreateUserGroup) (*UserGroup, error) {
	exists, err := c.checkExistUserGroup(action.Name)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, USER_GROUP_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newUserGroup := UserGroup{}
	err = json.Unmarshal(body, &newUserGroup)
	if err != nil {
		return nil, err
	}

	return &newUserGroup, nil
}

func (c *GFWClient) checkExistUserGroup(name string) (*UserGroup, error) {
	actions, err := c.GetUserGroups()
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

func (c *GFWClient) CreateUserGroupRole(userGroupRole CreateUserGroupRole) error {
	// obtain Role
	userGroup, err := c.GetUserGroup(strconv.Itoa(userGroupRole.UserGroupID))
	if err != nil {
		return err
	}
	var permsToCreate []int
	var permsToDelete []int
	for _, p := range userGroup.Roles {
		if !utils.Exists(p.ID, userGroupRole.Roles) {
			permsToDelete = append(permsToDelete, p.ID)
		}
	}

	for _, p := range userGroupRole.Roles {
		if !existsInRole(p, userGroup.Roles) {
			permsToCreate = append(permsToCreate, p)
		}
	}

	for _, createPermId := range permsToCreate {
		_, err := c.AddRoleInUserGroup(userGroupRole.UserGroupID, createPermId)
		if err != nil {
			return err
		}
	}
	for _, delPermId := range permsToDelete {
		_, err := c.DeleteRoleInUserGroup(userGroupRole.UserGroupID, delPermId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *GFWClient) DeleteUserGroupRole(userGroupId string) error {
	userGroup, err := c.GetUserGroup(userGroupId)
	if err != nil {
		return err
	}
	for _, r := range userGroup.Roles {
		_, err := c.DeleteRoleInUserGroup(userGroup.ID, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GFWClient) AddRoleInUserGroup(userGroupId, roleId int) (*Role, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%d/role/%d", c.HostURL, USER_GROUP_PATH, userGroupId, roleId), nil)
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

func (c *GFWClient) DeleteRoleInUserGroup(roleId, permissionId int) (*Role, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%d/role/%d", c.HostURL, USER_GROUP_PATH, roleId, permissionId), nil)
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

func existsInRole(i int, array []Role) bool {

	for _, v := range array {
		if i == v.ID {
			return true
		}
	}
	return false
}

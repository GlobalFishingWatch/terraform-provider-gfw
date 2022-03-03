package api

type Action struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type CreateAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Resource struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type CreateResource struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Permission struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"createdAt"`
	Resource    Resource `json:"resource"`
	Action      Action   `json:"action"`
}

type CreatePermission struct {
	Name        string `json:"name"`
	Action      int    `json:"actionId"`
	Resource    int    `json:"resourceId"`
	Description string `json:"description"`
}

type Role struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   string       `json:"createdAt"`
	Permissions []Permission `json:"permissions"`
}

type CreateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRolePermissions struct {
	RoleID      int
	Permissions []int
}

type UserGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	CreatedAt   string `json:"createdAt"`
	Roles       []Role `json:"roles"`
}

type CreateUserGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type CreateUserGroupRole struct {
	UserGroupID int
	Roles       []int
}

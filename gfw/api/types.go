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

type DatasetDocumentation struct {
	Type     string   `json:"type,omitempty"`
	Enable   bool     `json:"enable,omitempty"`
	Status   string   `json:"status,omitempty"`
	Queries  []string `json:"queries,omitempty"`
	Provider string   `json:"provider,omitempty"`
}
type DatasetConfigurationRange struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
}
type DatasetConfiguration struct {
	ApiSupportedVersions    []string                   `json:"apiSupportedVersions,omitempty"`
	InteractionColumns      []string                   `json:"interactionColumns,omitempty"`
	InteractionGroupColumns []string                   `json:"interactionGroupColumns,omitempty"`
	MaxZoom                 int                        `json:"maxZoom,omitempty"`
	Source                  string                     `json:"source,omitempty"`
	Function                string                     `json:"function,omitempty"`
	Type                    string                     `json:"type,omitempty"`
	GeometryColumn          string                     `json:"geometryColumn,omitempty"`
	DatabaseInstance        string                     `json:"databaseInstance,omitempty"`
	Project                 string                     `json:"project,omitempty"`
	Dataset                 string                     `json:"dataset,omitempty"`
	Table                   string                     `json:"table,omitempty"`
	Bucket                  string                     `json:"bucket,omitempty"`
	Folder                  string                     `json:"folder,omitempty"`
	Intervals               []string                   `json:"intervals,omitempty"`
	NumLayers               int                        `json:"numLayers,omitempty"`
	Index                   string                     `json:"index,omitempty"`
	Version                 int                        `json:"version,omitempty"`
	Translate               bool                       `json:"translate,omitempty"`
	Documentation           *DatasetDocumentation      `json:"documentation,omitempty"`
	Fields                  []string                   `json:"fields,omitempty"`
	GeometryType            string                     `json:"geometryType,omitempty"`
	PropertyToInclude       string                     `json:"propertyToInclude,omitempty"`
	PropertyToIncludeRange  *DatasetConfigurationRange `json:"propertyToIncludeRange,omitempty"`
	FilePath                string                     `json:"filePath,omitempty"`
	Srid                    string                     `json:"srid,omitempty"`
	Format                  string                     `json:"format,omitempty"`
	Latitude                string                     `json:"latitude,omitempty"`
	Longitude               string                     `json:"longitude,omitempty"`
	Timestamp               string                     `json:"timestamp,omitempty"`
}
type RelatedDataset struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Dataset struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	CreatedAt       string                  `json:"createdAt"`
	Type            string                  `json:"type"`
	Alias           []string                `json:"alias"`
	StartDate       string                  `json:"startDate"`
	EndDate         string                  `json:"endDate"`
	Unit            string                  `json:"unit"`
	Status          string                  `json:"status"`
	Category        string                  `json:"category"`
	Subcategory     string                  `json:"subcategory"`
	Source          string                  `json:"source"`
	Configuration   *DatasetConfiguration   `json:"configuration"`
	RelatedDatasets []RelatedDataset        `json:"relatedDatasets"`
	Schema          *map[string]interface{} `json:"schema"`
	FieldsAllowed   []string                `json:"fieldsAllowed"`
}

type CreateDataset struct {
	ID              string                  `json:"id,omitempty"`
	Name            string                  `json:"name,omitempty"`
	Description     string                  `json:"description,omitempty"`
	Type            string                  `json:"type,omitempty"`
	Alias           []string                `json:"alias,omitempty"`
	StartDate       string                  `json:"startDate,omitempty"`
	EndDate         string                  `json:"endDate,omitempty"`
	Unit            string                  `json:"unit,omitempty"`
	Status          string                  `json:"status,omitempty"`
	Category        string                  `json:"category,omitempty"`
	Subcategory     string                  `json:"subcategory,omitempty"`
	Source          string                  `json:"source,omitempty"`
	Configuration   *DatasetConfiguration   `json:"configuration,omitempty"`
	RelatedDatasets []RelatedDataset        `json:"relatedDatasets,omitempty"`
	Schema          *map[string]interface{} `json:"schema,omitempty"`
	FieldsAllowed   []string                `json:"fieldsAllowed,omitempty"`
}

type DataviewConfiguration struct {
	Type      string   `json:"type"`
	Color     string   `json:"color"`
	Datasets  []string `json:"datasets"`
	ColorRamp string   `json:"colorRamp"`
}

type Dataview struct {
	ID             int                       `json:"id"`
	Name           string                    `json:"name"`
	Slug           string                    `json:"slug"`
	Description    string                    `json:"description"`
	Category       string                    `json:"category"`
	App            string                    `json:"app"`
	CreatedAt      string                    `json:"createdAt"`
	UpdatedAt      string                    `json:"updatedAt"`
	Config         *DataviewConfiguration    `json:"config"`
	InfoConfig     *map[string]interface{}   `json:"infoConfig"`
	DatasetsConfig *[]map[string]interface{} `json:"datasetsConfig"`
}

type CreateDataview struct {
	Name           string                    `json:"name,omitempty"`
	Slug           string                    `json:"slug,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Category       string                    `json:"category,omitempty"`
	App            string                    `json:"app,omitempty"`
	CreatedAt      string                    `json:"createdAt,omitempty"`
	UpdatedAt      string                    `json:"updatedAt,omitempty"`
	Config         *DataviewConfiguration    `json:"config,omitempty"`
	InfoConfig     *map[string]interface{}   `json:"infoConfig,omitempty"`
	DatasetsConfig *[]map[string]interface{} `json:"datasetsConfig,omitempty"`
}

type WorkspaceViewport struct {
	Zoom      float64 `json:"zoom"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WorkspaceDataviewInstance struct {
	ID             string                   `json:"id"`
	Config         *map[string]interface{}  `json:"config"`
	DataviewID     int                      `json:"dataviewId"`
	DatasetsConfig []map[string]interface{} `json:"datasetsConfig"`
}

type Workspace struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Category    string                  `json:"category"`
	App         string                  `json:"app"`
	Aoi         string                  `json:"aoi"`
	StartAt     string                  `json:"startAt"`
	EndAt       string                  `json:"endAt"`
	Public      bool                    `json:"public"`
	Viewport    *WorkspaceViewport      `json:"viewport"`
	State       *map[string]interface{} `json:"state"`
	Dataviews   []struct {
		ID int `json:"id"`
	} `json:"dataviews"`
	DataviewInstances *[]WorkspaceDataviewInstance `json:"dataviewInstances"`
	CreatedAt         string                       `json:"createdAt"`
}

type CreateWorkspace struct {
	Name              string                       `json:"name,omitempty"`
	Description       string                       `json:"description,omitempty"`
	Category          string                       `json:"category,omitempty"`
	App               string                       `json:"app,omitempty"`
	Aoi               string                       `json:"aoi,omitempty"`
	StartAt           string                       `json:"startAt,omitempty"`
	EndAt             string                       `json:"endAt,omitempty"`
	Public            bool                         `json:"public,omitempty"`
	Viewport          *WorkspaceViewport           `json:"viewport,omitempty"`
	State             *map[string]interface{}      `json:"state,omitempty"`
	Dataviews         []int                        `json:"dataviews,omitempty"`
	DataviewInstances *[]WorkspaceDataviewInstance `json:"dataviewInstances,omitempty"`
}

type Pagination[T any] struct {
	Total      int                    `json:"total"`
	Limit      *int                   `json:"limit"`
	Offset     *int                   `json:"offset"`
	NextOffset *int                   `json:"nextOffset"`
	Metadata   map[string]interface{} `json:"metadata"`
	Entries    []T                    `json:"entries"`
}

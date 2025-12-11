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

type InsightSources struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Insight string `json:"insight,omitempty"`
}

type DatasetConfigurationRange struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
}

type DOIConfiguration struct {
	DOI        string `json:"doi,omitempty"`
	ConceptDOI int    `json:"conceptDOI,omitempty"`
}

// Filter Configuration
type FilterConfig struct {
	Label           string   `json:"label,omitempty"`
	ID              string   `json:"id,omitempty"`
	Type            string   `json:"type,omitempty"`
	Required        bool     `json:"required,omitempty"`
	Array           bool     `json:"array,omitempty"`
	Enum            []string `json:"enum,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	Format          string   `json:"format,omitempty"`
	MaxLength       int      `json:"maxLength,omitempty"`
	MinLength       int      `json:"minLength,omitempty"`
	Max             float64  `json:"max,omitempty"`
	Min             float64  `json:"min,omitempty"`
	SingleSelection bool     `json:"singleSelection,omitempty"`
	Operation       string   `json:"operation,omitempty"`
}

type DatasetFilters struct {
	Fourwings []FilterConfig `json:"4wings,omitempty"`
	Events    []FilterConfig `json:"events,omitempty"`
	Vessels   []FilterConfig `json:"vessels,omitempty"`
	Tracks    []FilterConfig `json:"tracks,omitempty"`
}

// Context Layer Configuration
type ContextLayerV1Config struct {
	ImportLogs string   `json:"importLogs,omitempty"`
	Srid       string   `json:"srid,omitempty"`
	Format     string   `json:"format,omitempty"`
	Fields     []string `json:"fields,omitempty"`
	FilePath   string   `json:"filePath,omitempty"`
	IDProperty string   `json:"idProperty,omitempty"`
}

// User Context Layer Configuration
type UserContextLayerV1Config struct {
	Table           string   `json:"table,omitempty"`
	ImportLogs      string   `json:"importLogs,omitempty"`
	Srid            string   `json:"srid,omitempty"`
	Format          string   `json:"format,omitempty"`
	Fields          []string `json:"fields,omitempty"`
	FilePath        string   `json:"filePath,omitempty"`
	IDProperty      string   `json:"idProperty,omitempty"`
	ValuePropertyID string   `json:"valuePropertyId,omitempty"`
}

// Temporal Context Layer Configuration
type TemporalContextLayerV1Config struct {
	Dataset string `json:"dataset,omitempty"`
	Project string `json:"project,omitempty"`
	Source  string `json:"source,omitempty"`
	Table   string `json:"table,omitempty"`
}

// User Tracks Configuration
type UserTracksV1Config struct {
	FilePath   string `json:"filePath,omitempty"`
	IDProperty string `json:"idProperty,omitempty"`
}

// PM Tiles Configuration
type PmTilesV1Config struct {
	FilePath   string `json:"filePath,omitempty"`
	IDProperty string `json:"idProperty,omitempty"`
}

// Extra Property Position Tiles
type ExtraPropertyPositionTiles struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// Events Configuration
type EventsV1Config struct {
	Table    string `json:"table,omitempty"`
	Dataset  string `json:"dataset,omitempty"`
	Project  string `json:"project,omitempty"`
	Function string `json:"function,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	MaxZoom  int    `json:"maxZoom,omitempty"`
	Source   string `json:"source,omitempty"`
}

// 4wings Configuration
type FourwingsV1Config struct {
	ExtraPropertiesPositionTiles []ExtraPropertyPositionTiles `json:"extraPropertiesPositionTiles,omitempty"`
	ReportGroupings              []string                     `json:"reportGroupings,omitempty"`
	Table                        string                       `json:"table,omitempty"`
	Dataset                      string                       `json:"dataset,omitempty"`
	MaxZoom                      int                          `json:"maxZoom,omitempty"`
	Project                      string                       `json:"project,omitempty"`
	Function                     string                       `json:"function,omitempty"`
	Intervals                    []string                     `json:"intervals,omitempty"`
	TTL                          int                          `json:"ttl,omitempty"`
	Max                          float64                      `json:"max,omitempty"`
	Min                          float64                      `json:"min,omitempty"`
	TileScale                    float64                      `json:"tileScale,omitempty"`
	TileOffset                   float64                      `json:"tileOffset,omitempty"`
	InternalScale                float64                      `json:"internalScale,omitempty"`
	InternalOffset               float64                      `json:"internalOffset,omitempty"`
	GeeBand                      string                       `json:"geeBand,omitempty"`
	GeeImages                    []string                     `json:"geeImages,omitempty"`
	InteractionColumns           []string                     `json:"interactionColumns,omitempty"`
	InteractionGroupColumns      []string                     `json:"interactionGroupColumns,omitempty"`
	TemporalAggregation          bool                         `json:"temporalAggregation,omitempty"`
	Source                       string                       `json:"source,omitempty"`
	Bucket                       string                       `json:"bucket,omitempty"`
	Folder                       string                       `json:"folder,omitempty"`
}

// Tracks Configuration
type TracksV1Config struct {
	Bucket           string `json:"bucket,omitempty"`
	Folder           string `json:"folder,omitempty"`
	DatabaseInstance string `json:"databaseInstance,omitempty"`
	Table            string `json:"table,omitempty"`
}

// Frontend Configuration
type FrontendConfig struct {
	MaxZoom            int      `json:"maxZoom,omitempty"`
	Translate          bool     `json:"translate,omitempty"`
	Max                float64  `json:"max,omitempty"`
	Min                float64  `json:"min,omitempty"`
	DisableInteraction bool     `json:"disableInteraction,omitempty"`
	Latitude           string   `json:"latitude,omitempty"`
	Longitude          string   `json:"longitude,omitempty"`
	StartTime          string   `json:"startTime,omitempty"`
	EndTime            string   `json:"endTime,omitempty"`
	Timestamp          string   `json:"timestamp,omitempty"`
	GeometryType       string   `json:"geometryType,omitempty"`
	SourceFormat       string   `json:"sourceFormat,omitempty"`
	TimeFilterType     string   `json:"timeFilterType,omitempty"`
	ValueProperties    []string `json:"valueProperties,omitempty"`
	PolygonColor       string   `json:"polygonColor,omitempty"`
	PointSize          string   `json:"pointSize,omitempty"`
	LineID             string   `json:"lineId,omitempty"`
	SegmentID          string   `json:"segmentId,omitempty"`
}

// Vessels Configuration
type VesselsV1Config struct {
	Index      string  `json:"index,omitempty"`
	IndexBoost float64 `json:"indexBoost,omitempty"`
	Table      string  `json:"table,omitempty"`
}

// Insights Configuration
type InsightsV1Config struct {
	Sources []InsightSources `json:"sources,omitempty"`
}

// Bulk Download Configuration
type BulkDownloadV1Config struct {
	GcsUri            string `json:"gcsUri,omitempty"`
	Path              string `json:"path,omitempty"`
	Format            string `json:"format,omitempty"`
	Compressed        bool   `json:"compressed,omitempty"`
	LatitudeProperty  string `json:"latitudeProperty,omitempty"`
	LongitudeProperty string `json:"longitudeProperty,omitempty"`
}

// Data Download Configuration
type DataDownloadV1Config struct {
	EmailGroups []string `json:"emailGroups,omitempty"`
	GcsFolder   string   `json:"gcsFolder,omitempty"`
	Doi         string   `json:"doi,omitempty"`
	ConceptDOI  int      `json:"conceptDOI,omitempty"`
}

// Thumbnails Configuration
type ThumbnailsV1Config struct {
	Extensions []string `json:"extensions,omitempty"`
	Bucket     string   `json:"bucket,omitempty"`
	Folder     string   `json:"folder,omitempty"`
	Scale      float64  `json:"scale,omitempty"`
}

type DatasetConfiguration struct {
	ApiSupportedVersions   []string                      `json:"apiSupportedVersions,omitempty"`
	ContextLayerV1         *ContextLayerV1Config         `json:"contextLayerV1,omitempty"`
	UserContextLayerV1     *UserContextLayerV1Config     `json:"userContextLayerV1,omitempty"`
	TemporalContextLayerV1 *TemporalContextLayerV1Config `json:"temporalContextLayerV1,omitempty"`
	UserTracksV1           *UserTracksV1Config           `json:"userTracksV1,omitempty"`
	PmTilesV1              *PmTilesV1Config              `json:"pmTilesV1,omitempty"`
	EventsV1               *EventsV1Config               `json:"eventsV1,omitempty"`
	FourwingsV1            *FourwingsV1Config            `json:"fourwingsV1,omitempty"`
	TracksV1               *TracksV1Config               `json:"tracksV1,omitempty"`
	Frontend               *FrontendConfig               `json:"frontend,omitempty"`
	VesselsV1              *VesselsV1Config              `json:"vesselsV1,omitempty"`
	InsightsV1             *InsightsV1Config             `json:"insightsV1,omitempty"`
	BulkDownloadV1         *BulkDownloadV1Config         `json:"bulkDownloadV1,omitempty"`
	DataDownloadV1         *DataDownloadV1Config         `json:"dataDownloadV1,omitempty"`
	ThumbnailsV1           *ThumbnailsV1Config           `json:"thumbnailsV1,omitempty"`
}
type RelatedDataset struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Dataset struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	CreatedAt       string                `json:"createdAt"`
	Type            string                `json:"type"`
	Alias           []string              `json:"alias"`
	StartDate       string                `json:"startDate"`
	EndDate         string                `json:"endDate"`
	Unit            string                `json:"unit"`
	Status          string                `json:"status"`
	Category        string                `json:"category"`
	Subcategory     string                `json:"subcategory"`
	Source          string                `json:"source"`
	Configuration   *DatasetConfiguration `json:"configuration"`
	RelatedDatasets []RelatedDataset      `json:"relatedDatasets"`
	Filters         *DatasetFilters       `json:"filters,omitempty"`
	Documentation   *DatasetDocumentation `json:"documentation,omitempty"`
}

type CreateDataset struct {
	ID              string                `json:"id,omitempty"`
	Name            string                `json:"name,omitempty"`
	Description     string                `json:"description,omitempty"`
	Type            string                `json:"type,omitempty"`
	Alias           []string              `json:"alias,omitempty"`
	StartDate       string                `json:"startDate,omitempty"`
	EndDate         string                `json:"endDate,omitempty"`
	Unit            string                `json:"unit,omitempty"`
	Status          string                `json:"status,omitempty"`
	Category        string                `json:"category,omitempty"`
	Subcategory     string                `json:"subcategory,omitempty"`
	Source          string                `json:"source,omitempty"`
	Configuration   *DatasetConfiguration `json:"configuration,omitempty"`
	RelatedDatasets []RelatedDataset      `json:"relatedDatasets,omitempty"`
	Filters         *DatasetFilters       `json:"filters,omitempty"`
	Documentation   *DatasetDocumentation `json:"documentation,omitempty"`
}

type DataviewLayer struct {
	ID      string `json:"id"`
	Dataset string `json:"dataset"`
}

type DataviewConfiguration struct {
	Type                 string                  `json:"type,omitempty"`
	Color                string                  `json:"color,omitempty"`
	Datasets             []string                `json:"datasets,omitempty"`
	ColorRamp            string                  `json:"colorRamp,omitempty"`
	Filters              *map[string]interface{} `json:"filters,omitempty"`
	ClusterMaxZoomLevels *map[string]interface{} `json:"clusterMaxZoomLevels,omitempty"`
	Pickable             bool                    `json:"pickable,omitempty"`
	MaxZoom              int                     `json:"maxZoom,omitempty"`
	AggregationOperation string                  `json:"aggregationOperation,omitempty"`
	Layers               []DataviewLayer         `json:"layers,omitempty"`
	Breaks               []float64               `json:"breaks,omitempty"`
	Intervals            []string                `json:"intervals,omitempty"`
}

type Dataview struct {
	ID             int                       `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Slug           string                    `json:"slug,omitempty"`
	Description    string                    `json:"description,omitempty"`
	Category       string                    `json:"category,omitempty"`
	App            string                    `json:"app,omitempty"`
	CreatedAt      string                    `json:"createdAt,omitempty"`
	UpdatedAt      string                    `json:"updatedAt,omitempty"`
	Config         *DataviewConfiguration    `json:"config,omitempty"`
	InfoConfig     *map[string]interface{}   `json:"infoConfig,omitempty"`
	EventsConfig   *map[string]interface{}   `json:"eventsConfig,omitempty"`
	FiltersConfig  *map[string]interface{}   `json:"filtersConfig,omitempty"`
	DatasetsConfig *[]map[string]interface{} `json:"datasetsConfig,omitempty"`
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
	EventsConfig   *map[string]interface{}   `json:"eventsConfig,omitempty"`
	FiltersConfig  *map[string]interface{}   `json:"filtersConfig,omitempty"`
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
	Category       string                   `json:"category"`
	DataviewID     string                   `json:"dataviewId"`
	DatasetsConfig []map[string]interface{} `json:"datasetsConfig"`
}

type Workspace struct {
	ID                string                       `json:"id"`
	Name              string                       `json:"name"`
	Description       string                       `json:"description"`
	Category          string                       `json:"category"`
	App               string                       `json:"app"`
	Aoi               string                       `json:"aoi"`
	StartAt           string                       `json:"startAt"`
	EndAt             string                       `json:"endAt"`
	Public            bool                         `json:"public"`
	Viewport          *WorkspaceViewport           `json:"viewport"`
	State             *map[string]interface{}      `json:"state"`
	DataviewInstances *[]WorkspaceDataviewInstance `json:"dataviewInstances"`
	CreatedAt         string                       `json:"createdAt"`
}

type CreateWorkspace struct {
	ID                string                       `json:"id,omitempty"`
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

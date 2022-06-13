package gfw

import (
	"context"
	"encoding/json"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var DATASET_TYPES []string = []string{
	"tracks:v1",
	"vessels:v1",
	"events:v1",
	"ports:v1",
	"4wings:v1",
	"user-tracks:v1",
	"user-context-layer:v1",
	"data-download:v1",
	"temporal-context-layer:v1",
}
var DATASET_CATEGORIES []string = []string{
	"activity",
	"context",
	"context-layer",
	"detection",
	"environment",
	"event",
	"vessel",
}
var DATASET_SUBCATEGORIES []string = []string{
	"track",
	"animal",
	"loitering",
	"presence",
	"port_visit",
	"fishing",
	"info",
	"viirs",
	"sar",
	"encounter",
	"salinity",
	"chlorophyl",
	"water-temperature",
	"user",
}
var DATASET_UNITS []string = []string{
	"unit",
	"TBD",
	"probability",
	"hours",
	"days",
	"mg/m^3",
	"PSU",
	"ºC",
	"detections",
	"habitat suitability",
	"NA",
}
var DATASET_CONFIGURATION_GEOMETRY_TYPES []string = []string{"tracks", "polygons", "points"}
var DATASET_CONFIGURATION_FORMATS []string = []string{"geojson"}

func resourceDataset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasetCreate,
		ReadContext:   resourceDatasetRead,
		UpdateContext: resourceDatasetUpdate,
		DeleteContext: resourceDatasetDelete,
		Schema: map[string]*schema.Schema{
			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_TYPES, false),
			},
			"alias": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type:     schema.TypeList,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(DATASET_UNITS, false),
			},
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_CATEGORIES, false),
			},
			"subcategory": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATASET_SUBCATEGORIES, false),
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related_datasets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(DATASET_TYPES, false),
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_supported_versions": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"interaction_columns": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"interaction_group_columns": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"max_zoom": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  12,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"geometry_column": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"database_instance": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dataset": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"table": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"folder": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"intervals": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"num_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"num_layers": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"index": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"translate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"documentation": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"queries": {
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Type:     schema.TypeList,
										Optional: true,
									},
									"provider": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"fields": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"geometry_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(DATASET_CONFIGURATION_GEOMETRY_TYPES, false),
						},
						"property_to_include": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"property_to_include_range": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"max": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"srid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(DATASET_CONFIGURATION_FORMATS, false),
						},
						"latitude": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"longitude": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"fields_allowed": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type:     schema.TypeList,
				Optional: true,
			},
			"schema": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceDatasetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	id := d.Get("dataset_id").(string)
	dataset, err := schemaToDataset(d)
	if err != nil {
		return diag.FromErr(err)
	}
	dataset.ID = id
	datasetCreated, err := c.CreateDataset(dataset)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(datasetCreated.ID)
	resourceDatasetRead(ctx, d, m)
	return diags
}

func resourceDatasetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	datasetId := d.Id()
	c := m.(*api.GFWClient)
	dataset, err := c.GetDataset(datasetId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", dataset.Name)
	d.Set("description", dataset.Description)
	d.Set("created_at", dataset.CreatedAt)
	d.Set("type", dataset.Type)
	d.Set("alias", dataset.Alias)
	d.Set("start_date", dataset.StartDate)
	d.Set("end_date", dataset.EndDate)
	d.Set("unit", dataset.Unit)
	d.Set("category", dataset.Category)
	d.Set("subcategory", dataset.Subcategory)
	d.Set("source", dataset.Source)
	d.Set("fields_allowed", dataset.FieldsAllowed)
	d.Set("type", dataset.Type)

	if dataset.Schema != nil {
		jsonStr, err := json.Marshal(dataset.Schema)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("schema", string(jsonStr)); err != nil {
			return diag.FromErr(err)
		}
	}

	if dataset.Configuration != nil {
		configuration := flattenDatasetConfiguration(*dataset.Configuration)
		if err := d.Set("configuration", []interface{}{configuration}); err != nil {
			return diag.FromErr(err)
		}
	}

	relatedDatasets := flattenRelatedDatasets(dataset.RelatedDatasets)
	if err := d.Set("related_datasets", relatedDatasets); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatasetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dataset, err := schemaToDataset(d)
	if err != nil {
		return diag.FromErr(err)
	}
	datasetId := d.Id()
	c := m.(*api.GFWClient)
	err = c.UpdateDataset(datasetId, dataset)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDatasetRead(ctx, d, m)
}

func schemaToDataset(d *schema.ResourceData) (api.CreateDataset, error) {
	dataset := api.CreateDataset{}
	if d.HasChange("name") {
		dataset.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		dataset.Description = d.Get("description").(string)
	}
	if d.HasChange("type") {
		dataset.Type = d.Get("type").(string)
	}
	if d.HasChange("unit") {
		dataset.Unit = d.Get("unit").(string)
	}
	if d.HasChange("category") {
		dataset.Category = d.Get("category").(string)
	}
	if d.HasChange("subcategory") {
		dataset.Subcategory = d.Get("subcategory").(string)
	}
	if d.HasChange("source") {
		dataset.Source = d.Get("source").(string)
	}
	if d.HasChange("start_date") {
		dataset.StartDate = d.Get("start_date").(string)
	}
	if d.HasChange("end_date") {
		dataset.EndDate = d.Get("end_date").(string)
	}
	if d.HasChange("alias") && d.Get("alias") != nil {
		dataset.Alias = utils.ConvertArrayInterfaceToArrayString(d.Get("alias").([]interface{}))
	}
	if d.HasChange("fields_allowed") && d.Get("fields_allowed") != nil {
		dataset.FieldsAllowed = utils.ConvertArrayInterfaceToArrayString(d.Get("fields_allowed").([]interface{}))
	}
	if d.HasChange("schema") && d.Get("schema") != nil {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("schema").(string)), &obj)
		if err != nil {
			return api.CreateDataset{}, err
		}
		dataset.Schema = &obj
	}
	if d.HasChange("configuration") && d.Get("configuration") != nil {
		configuration := d.Get("configuration").([]interface{})
		if len(configuration) > 0 {
			config := schemaToDatasetConfiguration(configuration[0].(map[string]interface{}))
			dataset.Configuration = &config
		}
	}
	if d.HasChange("related_datasets") && d.Get("related_datasets") != nil {
		relatedDatasets := schemaToRelatedDatasets(d.Get("related_datasets").([]interface{}))
		dataset.RelatedDatasets = relatedDatasets
	}

	return dataset, nil
}

func schemaToDatasetConfiguration(schema map[string]interface{}) api.DatasetConfiguration {
	config := api.DatasetConfiguration{
		Source:            schema["source"].(string),
		Function:          schema["function"].(string),
		Type:              schema["type"].(string),
		GeometryColumn:    schema["geometry_column"].(string),
		DatabaseInstance:  schema["database_instance"].(string),
		Project:           schema["project"].(string),
		Dataset:           schema["dataset"].(string),
		Table:             schema["table"].(string),
		Bucket:            schema["bucket"].(string),
		Folder:            schema["folder"].(string),
		Index:             schema["index"].(string),
		GeometryType:      schema["geometry_type"].(string),
		PropertyToInclude: schema["property_to_include"].(string),
		FilePath:          schema["file_path"].(string),
		Srid:              schema["srid"].(string),
		Format:            schema["format"].(string),
		Latitude:          schema["latitude"].(string),
		Longitude:         schema["longitude"].(string),
		Timestamp:         schema["timestamp"].(string),
		NumBytes:          schema["num_bytes"].(int),
	}
	if val, ok := schema["max_zoom"]; ok {
		maxZoom := val.(int)
		config.MaxZoom = maxZoom
	}
	if val, ok := schema["translate"]; ok {
		translate := val.(bool)
		config.Translate = translate
	}
	if val, ok := schema["num_layers"]; ok {
		numLayers := val.(int)
		config.NumLayers = numLayers
	}
	if val, ok := schema["version"]; ok {
		version := val.(int)
		config.Version = version
	}
	if val, ok := schema["api_supported_versions"]; ok {
		config.ApiSupportedVersions = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["intervals"]; ok {
		config.Intervals = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["interaction_group_columns"]; ok {
		config.InteractionGroupColumns = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["interaction_columns"]; ok {
		config.InteractionColumns = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["fields"]; ok {
		config.Fields = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	if val, ok := schema["documentation"]; ok {
		documentationArray := val.([]interface{})
		if len(documentationArray) > 0 {
			doc := schemaToDatasetDocumentation(documentationArray[0].(map[string]interface{}))
			config.Documentation = &doc
		}
	}
	if val, ok := schema["property_to_include_range"]; ok {
		propertyToIncludeRangeArray := val.([]interface{})
		if len(propertyToIncludeRangeArray) > 0 {
			prop := schemaToDatasetConfigurationRange(propertyToIncludeRangeArray[0].(map[string]interface{}))
			config.PropertyToIncludeRange = &prop
		}
	}

	return config
}

func schemaToDatasetConfigurationRange(schema map[string]interface{}) api.DatasetConfigurationRange {
	doc := api.DatasetConfigurationRange{
		Min: schema["type"].(float64),
		Max: schema["enable"].(float64),
	}
	return doc
}
func schemaToDatasetDocumentation(schema map[string]interface{}) api.DatasetDocumentation {
	doc := api.DatasetDocumentation{
		Type:     schema["type"].(string),
		Status:   schema["status"].(string),
		Provider: schema["provider"].(string),
	}
	if val, ok := schema["enable"]; ok {
		enable := val.(bool)
		doc.Enable = enable
	}
	if val, ok := schema["queries"]; ok {
		doc.Queries = utils.ConvertArrayInterfaceToArrayString(val.([]interface{}))
	}
	return doc
}

func schemaToRelatedDatasets(schema []interface{}) []api.RelatedDataset {
	relatedDatasets := make([]api.RelatedDataset, len(schema))
	for i, s := range schema {

		relatedDatasets[i] = api.RelatedDataset{
			ID:   s.(map[string]interface{})["id"].(string),
			Type: s.(map[string]interface{})["type"].(string),
		}
	}
	return relatedDatasets
}

func resourceDatasetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	datasetId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteDataset(datasetId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenDatasetConfiguration(config api.DatasetConfiguration) interface{} {
	a := make(map[string]interface{})

	a["api_supported_versions"] = config.ApiSupportedVersions
	a["interaction_columns"] = config.InteractionColumns
	a["interaction_group_columns"] = config.InteractionGroupColumns
	a["max_zoom"] = config.MaxZoom
	a["source"] = config.Source
	a["function"] = config.Function
	a["type"] = config.Type
	a["geometry_column"] = config.GeometryColumn
	a["database_instance"] = config.DatabaseInstance
	a["project"] = config.Project
	a["dataset"] = config.Dataset
	a["table"] = config.Table
	a["bucket"] = config.Bucket
	a["folder"] = config.Folder
	a["intervals"] = config.Intervals
	a["num_layers"] = config.NumLayers
	a["index"] = config.Index
	a["version"] = config.Version
	a["translate"] = config.Translate
	a["num_bytes"] = config.NumBytes
	if config.Documentation != nil {
		a["documentation"] = []interface{}{flattenDatasetDocumentation(*config.Documentation)}
	}
	a["fields"] = config.Fields
	a["geometry_type"] = config.GeometryType
	a["property_to_include"] = config.PropertyToInclude
	if config.PropertyToIncludeRange != nil {
		a["property_to_include_range"] = []interface{}{flattenDatasetConfigurationRange(*config.PropertyToIncludeRange)}
	}
	a["file_path"] = config.FilePath
	a["srid"] = config.Srid
	a["format"] = config.Format
	a["latitude"] = config.Latitude
	a["longitude"] = config.Longitude
	a["timestamp"] = config.Timestamp

	return a
}

func flattenRelatedDatasets(relatedDatasets []api.RelatedDataset) []map[string]interface{} {
	list := make([]map[string]interface{}, len(relatedDatasets))
	for i, rd := range relatedDatasets {
		list[i] = map[string]interface{}{}
		list[i]["id"] = rd.ID
		list[i]["type"] = rd.Type

	}

	return list
}

func flattenDatasetDocumentation(doc api.DatasetDocumentation) interface{} {
	a := make(map[string]interface{})

	a["type"] = doc.Type
	a["enable"] = doc.Enable
	a["status"] = doc.Status
	a["queries"] = doc.Queries
	a["provider"] = doc.Provider

	return a
}

func flattenDatasetConfigurationRange(r api.DatasetConfigurationRange) interface{} {
	a := make(map[string]interface{})

	a["max"] = r.Max
	a["min"] = r.Min

	return a
}

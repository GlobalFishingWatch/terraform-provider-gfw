package gfw

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var DATAVIEW_TYPES []string = []string{"context", "environment", "activity", "detections", "events", "vessels"}
var DATAVIEW_APPS []string = []string{"fishing-map", "vessel-history"}
var DATAVIEW_CONFIG_INTERVALS []string = []string{"hours", "day", "10days", "month"}
var DATAVIEW_CONFIG_TYPES []string = []string{
	"BASEMAP",
	"HEATMAP",
	"HEATMAP_ANIMATED",
	"TRACK",
	"CONTEXT",
	"USER_CONTEXT",
	"TILE_CLUSTER",
	"BACKGROUND",
	"POLYGONS",
	"USER_CONTEXT",
	"USER_POINTS",
	"VESSEL_EVENTS",
	"VESSEL_EVENTS_SHAPES",
}

func resourceDataview() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataviewCreate,
		ReadContext:   resourceDataviewRead,
		UpdateContext: resourceDataviewUpdate,
		DeleteContext: resourceDataviewDelete,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(DATAVIEW_TYPES, false),
			},
			"app": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(DATAVIEW_APPS, false),
			},

			"datasets_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
				},
			},
			"info_config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"events_config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"filters_config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datasets": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"intervals": {
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(DATAVIEW_CONFIG_INTERVALS, false),
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"breaks": {
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
							Type:     schema.TypeList,
							Optional: true,
						},
						"layers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"dataset": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(DATAVIEW_CONFIG_TYPES, false),
						},
						"aggregation_operation": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_zoom": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"color": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"color_ramp": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceDataviewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	dataview, err := schemaToDataview(d)
	if err != nil {
		return diag.FromErr(err)
	}
	dataviewCreated, err := c.CreateDataview(dataview)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(dataviewCreated.ID))
	resourceDataviewRead(ctx, d, m)
	return diags
}

func resourceDataviewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	dataviewId := d.Id()
	c := m.(*api.GFWClient)
	dataview, err := c.GetDataview(dataviewId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", dataview.Name)
	d.Set("description", dataview.Description)
	d.Set("created_at", dataview.CreatedAt)
	d.Set("slug", dataview.Slug)
	d.Set("app", dataview.App)
	d.Set("updated_at", dataview.UpdatedAt)
	d.Set("category", dataview.Category)

	if dataview.Config != nil {
		configuration := flattenDataviewConfiguration(*dataview.Config)
		if err := d.Set("config", []interface{}{configuration}); err != nil {
			return diag.FromErr(err)
		}
	}
	if dataview.InfoConfig != nil {
		jsonStr, err := json.Marshal(dataview.InfoConfig)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("info_config", string(jsonStr)); err != nil {
			return diag.FromErr(err)
		}
	}
	if dataview.FiltersConfig != nil {
		jsonStr, err := json.Marshal(dataview.FiltersConfig)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("filters_config", string(jsonStr)); err != nil {
			return diag.FromErr(err)
		}
	}
	if dataview.EventsConfig != nil {
		jsonStr, err := json.Marshal(dataview.EventsConfig)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("events_config", string(jsonStr)); err != nil {
			return diag.FromErr(err)
		}
	}
	if dataview.DatasetsConfig != nil {
		jsonStrArr := make([]string, len(*dataview.DatasetsConfig))
		for i, m := range *dataview.DatasetsConfig {
			jsonStr, err := json.Marshal(m)
			if err != nil {
				return diag.FromErr(err)
			}
			jsonStrArr[i] = string(jsonStr)
		}

		if err := d.Set("datasets_config", jsonStrArr); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceDataviewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dataview, err := schemaToDataview(d)
	if err != nil {
		return diag.FromErr(err)
	}
	dataview.Slug = ""
	dataviewId := d.Id()
	c := m.(*api.GFWClient)
	err = c.UpdateDataview(dataviewId, dataview)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDataviewRead(ctx, d, m)
}

func schemaToDataview(d *schema.ResourceData) (api.CreateDataview, error) {
	dataview := api.CreateDataview{}
	if d.HasChange("name") {
		dataview.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		dataview.Description = d.Get("description").(string)
	}
	if d.HasChange("app") {
		dataview.App = d.Get("app").(string)
	}
	if d.HasChange("category") {
		dataview.Category = d.Get("category").(string)
	}
	if d.HasChange("slug") {
		dataview.Slug = d.Get("slug").(string)
	}

	if d.HasChange("config") && d.Get("config") != nil {
		configuration := d.Get("config").([]interface{})
		if len(configuration) > 0 {
			config := schemaToDataviewConfiguration(configuration[0].(map[string]interface{}))
			dataview.Config = &config
		}
	}
	if d.HasChange("info_config") && d.Get("info_config") != nil {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("info_config").(string)), &obj)
		if err != nil {
			return api.CreateDataview{}, err
		}
		dataview.InfoConfig = &obj
	}
	if d.HasChange("filters_config") && d.Get("filters_config") != nil {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("filters_config").(string)), &obj)
		if err != nil {
			return api.CreateDataview{}, err
		}
		dataview.FiltersConfig = &obj
	}
	if d.HasChange("events_config") && d.Get("events_config") != nil {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("events_config").(string)), &obj)
		if err != nil {
			return api.CreateDataview{}, err
		}
		dataview.EventsConfig = &obj
	}
	if d.HasChange("datasets_config") && d.Get("datasets_config") != nil {
		list := d.Get("datasets_config").([]interface{})
		datasetsConfig := make([]map[string]interface{}, len(list))
		for i, m := range list {
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(m.(string)), &obj)
			if err != nil {
				return api.CreateDataview{}, err
			}
			datasetsConfig[i] = obj
		}

		dataview.DatasetsConfig = &datasetsConfig
	}

	return dataview, nil
}

func schemaToDataviewConfiguration(schema map[string]interface{}) api.DataviewConfiguration {
	config := api.DataviewConfiguration{
		Type:                 schema["type"].(string),
		Color:                schema["color"].(string),
		ColorRamp:            schema["color_ramp"].(string),
		Datasets:             utils.ConvertArrayInterfaceToArrayString(schema["datasets"].([]interface{})),
		Intervals:            utils.ConvertArrayInterfaceToArrayString(schema["intervals"].([]interface{})),
		Breaks:               utils.ConvertArrayInterfaceToArrayFloat(schema["breaks"].([]interface{})),
		AggregationOperation: schema["aggregation_operation"].(string),
	}
	if val, ok := schema["max_zoom"]; ok {
		config.MaxZoom = val.(int)
	}
	if val, ok := schema["layers"]; ok {
		layers := val.([]interface{})
		if len(layers) > 0 {
			arr := make([]api.DataviewLayer, len(layers))
			for i, l := range layers {
				arr[i] = schemaToDataviewLayer(l.(map[string]interface{}))
			}

			config.Layers = arr
		}
	}

	return config
}

func schemaToDataviewLayer(schema map[string]interface{}) api.DataviewLayer {
	el := api.DataviewLayer{
		ID: schema["id"].(string),
	}
	if val, ok := schema["dataset"]; ok {
		el.Dataset = val.(string)
	}
	return el
}

func resourceDataviewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	dataviewId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteDataview(dataviewId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenDataviewConfiguration(config api.DataviewConfiguration) interface{} {
	a := make(map[string]interface{})

	a["type"] = config.Type
	a["color"] = config.Color
	a["color_ramp"] = config.ColorRamp
	a["datasets"] = config.Datasets
	a["max_zoom"] = config.MaxZoom
	a["aggregation_operation"] = config.AggregationOperation
	a["breaks"] = config.Breaks
	a["intervals"] = config.Intervals

	if config.Layers != nil {
		layers := flattenDataviewLayer(config.Layers)
		a["layers"] = layers

	}
	return a
}

func flattenDataviewLayer(layers []api.DataviewLayer) interface{} {

	a := make([]map[string]interface{}, len(layers))
	for i, l := range layers {
		m := map[string]interface{}{}
		m["id"] = l.ID
		m["dataset"] = l.Dataset
		a[i] = m
	}

	return a
}

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

var DATAVIEW_TYPES []string = []string{"context", "environment", "fishing", "presence", "events", "vessels"}
var DATAVIEW_APPS []string = []string{"fishing-map", "vessel-history"}
var DATAVIEW_CONFIG_TYPES []string = []string{
	"BASEMAP",
	"HEATMAP",
	"HEATMAP_ANIMATED",
	"TRACK",
	"CONTEXT",
	"USER_CONTEXT",
	"TILE_CLUSTER",
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
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(DATAVIEW_CONFIG_TYPES, false),
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
		Type:      schema["type"].(string),
		Color:     schema["color"].(string),
		ColorRamp: schema["color_ramp"].(string),
		Datasets:  utils.ConvertArrayInterfaceToArrayString(schema["datasets"].([]interface{})),
	}

	return config
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

	return a
}

package gfw

import (
	"context"
	"encoding/json"
	"time"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var WORKSPACE_CATEGORIES []string = []string{"marine-reserves",
	"marine-manager",
	"fishing-activity",
	"country-portals"}

func resourceWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkspaceCreate,
		ReadContext:   resourceWorkspaceRead,
		UpdateContext: resourceWorkspaceUpdate,
		DeleteContext: resourceWorkspaceDelete,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(WORKSPACE_CATEGORIES, false),
			},
			"app": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"start_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"end_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.IsISOTime,
			},
			"aoi": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"viewport": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zoom": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"latitude": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"longitude": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					},
				},
			},
			"dataviews": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"dataview_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"dataview_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"datasets_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsJSON,
							},
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics
	id := d.Get("workspace_id").(string)
	workspace, err := schemaToWorkspace(d)
	if err != nil {
		return diag.FromErr(err)
	}
	workspace.ID = id
	workspaceCreated, err := c.CreateWorkspace(workspace)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(workspaceCreated.ID)
	resourceWorkspaceRead(ctx, d, m)
	return diags
}

func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	workspaceId := d.Id()
	c := m.(*api.GFWClient)
	workspace, err := c.GetWorkspace(workspaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", workspace.Name)
	d.Set("description", workspace.Description)
	d.Set("created_at", workspace.CreatedAt)
	d.Set("public", workspace.Public)
	d.Set("app", workspace.App)
	d.Set("start_at", workspace.StartAt)
	d.Set("end_at", workspace.EndAt)
	d.Set("category", workspace.Category)
	d.Set("aoi", workspace.Aoi)

	if workspace.Viewport != nil {
		configuration := flattenWorkspaceViewport(*workspace.Viewport)
		if err := d.Set("viewport", []interface{}{configuration}); err != nil {
			return diag.FromErr(err)
		}
	}
	if workspace.State != nil {
		jsonStr, err := json.Marshal(workspace.State)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("state", string(jsonStr)); err != nil {
			return diag.FromErr(err)
		}
	}
	if workspace.DataviewInstances != nil {
		dataviewInstances, err := flattenWorkspaceDataviewInstances(*workspace.DataviewInstances)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("dataview_instances", dataviewInstances); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workspace, err := schemaToWorkspace(d)
	if err != nil {
		return diag.FromErr(err)
	}
	workspaceId := d.Id()
	c := m.(*api.GFWClient)
	err = c.UpdateWorkspace(workspaceId, workspace)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceWorkspaceRead(ctx, d, m)
}

func schemaToWorkspace(d *schema.ResourceData) (api.CreateWorkspace, error) {
	workspace := api.CreateWorkspace{}
	if d.HasChange("name") {
		workspace.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		workspace.Description = d.Get("description").(string)
	}
	if d.HasChange("app") {
		workspace.App = d.Get("app").(string)
	}
	if d.HasChange("category") {
		workspace.Category = d.Get("category").(string)
	}
	if d.HasChange("public") {
		workspace.Public = d.Get("public").(bool)
	}
	if d.HasChange("start_at") {
		workspace.StartAt = d.Get("start_at").(string)
	}
	if d.HasChange("end_at") {
		workspace.EndAt = d.Get("end_at").(string)
	}
	if d.HasChange("aoi") {
		workspace.Aoi = d.Get("aoi").(string)
	}
	if d.HasChange("dataviews") {
		workspace.Dataviews = utils.ConvertArrayInterfaceToArrayInt(d.Get("dataviews").([]interface{}))
	}

	if d.HasChange("viewport") && d.Get("viewport") != nil {
		viewport := d.Get("viewport").([]interface{})
		if len(viewport) > 0 {
			viewportObj := schemaToWorkspaceViewport(viewport[0].(map[string]interface{}))
			workspace.Viewport = &viewportObj
		}
	}
	if d.HasChange("state") && d.Get("state") != nil {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("state").(string)), &obj)
		if err != nil {
			return api.CreateWorkspace{}, err
		}
		workspace.State = &obj
	}
	if d.HasChange("dataview_instances") && d.Get("dataview_instances") != nil {
		dataviewInstances := d.Get("dataview_instances").([]interface{})
		if len(dataviewInstances) > 0 {
			dataviewInstancesObj, err := schemaToWorkspaceDataviewInstances(dataviewInstances)
			if err != nil {
				return api.CreateWorkspace{}, err
			}
			workspace.DataviewInstances = &dataviewInstancesObj
		}
	}

	return workspace, nil
}

func schemaToWorkspaceViewport(schema map[string]interface{}) api.WorkspaceViewport {
	config := api.WorkspaceViewport{
		Zoom:      schema["zoom"].(float64),
		Latitude:  schema["latitude"].(float64),
		Longitude: schema["longitude"].(float64),
	}

	return config
}

func schemaToWorkspaceDataviewInstances(schema []interface{}) ([]api.WorkspaceDataviewInstance, error) {
	list := make([]api.WorkspaceDataviewInstance, len(schema))

	for i, inter := range schema {
		mp := inter.(map[string]interface{})
		list[i] = api.WorkspaceDataviewInstance{
			ID:         mp["id"].(string),
			Category:   mp["category"].(string),
			DataviewID: mp["dataview_id"].(string),
		}
		if mp["config"].(string) != "" {
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(mp["config"].(string)), &obj)
			if err != nil {
				return nil, err
			}
			list[i].Config = &obj
		}
		if mp["datasets_config"] != nil {
			listDatasets := mp["datasets_config"].([]interface{})
			datasetsConfig := make([]map[string]interface{}, len(listDatasets))
			for i, m := range listDatasets {
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(m.(string)), &obj)
				if err != nil {
					return nil, err
				}
				datasetsConfig[i] = obj
			}

			list[i].DatasetsConfig = datasetsConfig
		}

	}

	return list, nil
}

func resourceWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	workspaceId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteWorkspace(workspaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenWorkspaceViewport(config api.WorkspaceViewport) interface{} {
	a := make(map[string]interface{})

	a["zoom"] = config.Zoom
	a["latitude"] = config.Latitude
	a["longitude"] = config.Longitude

	return a
}
func flattenWorkspaceDataviewInstances(dataviewInstances []api.WorkspaceDataviewInstance) ([]interface{}, error) {
	list := make([]interface{}, len(dataviewInstances))
	for i, di := range dataviewInstances {

		a := make(map[string]interface{})
		a["id"] = di.ID
		a["category"] = di.Category
		a["dataview_id"] = di.DataviewID
		if di.Config != nil {
			jsonStr, err := json.Marshal(di.Config)
			if err != nil {
				return nil, err
			}
			a["config"] = string(jsonStr)
		}
		if di.DatasetsConfig != nil {
			jsonStrArr := make([]string, len(di.DatasetsConfig))
			for i, m := range di.DatasetsConfig {
				jsonStr, err := json.Marshal(m)
				if err != nil {
					return nil, err
				}
				jsonStrArr[i] = string(jsonStr)
			}
			a["datasets_config"] = jsonStrArr
		}
		list[i] = a

	}
	return list, nil
}

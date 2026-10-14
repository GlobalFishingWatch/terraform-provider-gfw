package gfw

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceReport() *schema.Resource {
	// The report workspace is stored inline in the report, but shares the gfw_workspace shape
	ws := resourceWorkspace().Schema
	return &schema.Resource{
		CreateContext: resourceReportCreate,
		ReadContext:   resourceReportRead,
		UpdateContext: resourceReportUpdate,
		DeleteContext: resourceReportDelete,
		Schema: map[string]*schema.Schema{
			"report_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"area_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dataset_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"workspace": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": ws["name"],
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"app":                ws["app"],
						"category":           ws["category"],
						"state":              ws["state"],
						"start_at":           ws["start_at"],
						"end_at":             ws["end_at"],
						"viewport":           ws["viewport"],
						"dataview_instances": ws["dataview_instances"],
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceReportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.GFWClient)
	workspace, err := schemaToReportWorkspace(d)
	if err != nil {
		return diag.FromErr(err)
	}
	report, err := c.CreateReport(api.CreateReport{
		ID:          d.Get("report_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		AreaID:      joinStringList(d.Get("area_ids")),
		DatasetID:   joinStringList(d.Get("dataset_ids")),
		Public:      d.Get("public").(bool),
		Workspace:   workspace,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(report.ID)
	return resourceReportRead(ctx, d, m)
}

func resourceReportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.GFWClient)
	report, err := c.GetReport(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", report.Name)
	d.Set("description", report.Description)
	d.Set("area_ids", strings.Split(report.AreaID, ","))
	d.Set("dataset_ids", strings.Split(report.DatasetID, ","))
	d.Set("public", report.Public)
	d.Set("created_at", report.CreatedAt)

	if report.Workspace != nil {
		workspace, err := flattenReportWorkspace(*report.Workspace)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("workspace", []interface{}{workspace}); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceReportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.GFWClient)
	// Always send the full report: the API replaces the workspace jsonb wholesale
	workspace, err := schemaToReportWorkspace(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = c.UpdateReport(d.Id(), api.UpdateReport{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		AreaID:      joinStringList(d.Get("area_ids")),
		DatasetID:   joinStringList(d.Get("dataset_ids")),
		Workspace:   workspace,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceReportRead(ctx, d, m)
}

func resourceReportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.GFWClient)
	if err := c.DeleteReport(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func joinStringList(v interface{}) string {
	return strings.Join(utils.ConvertArrayInterfaceToArrayString(v.([]interface{})), ",")
}

func schemaToReportWorkspace(d *schema.ResourceData) (*api.ReportWorkspace, error) {
	list := d.Get("workspace").([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}
	mp := list[0].(map[string]interface{})
	workspace := api.ReportWorkspace{
		ID:          mp["id"].(string),
		Name:        mp["name"].(string),
		Description: mp["description"].(string),
		App:         mp["app"].(string),
		Category:    mp["category"].(string),
		StartAt:     mp["start_at"].(string),
		EndAt:       mp["end_at"].(string),
	}
	if viewport := mp["viewport"].([]interface{}); len(viewport) > 0 {
		viewportObj := schemaToWorkspaceViewport(viewport[0].(map[string]interface{}))
		workspace.Viewport = &viewportObj
	}
	if state := mp["state"].(string); state != "" {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(state), &obj); err != nil {
			return nil, err
		}
		workspace.State = &obj
	}
	if dataviewInstances := mp["dataview_instances"].([]interface{}); len(dataviewInstances) > 0 {
		dataviewInstancesObj, err := schemaToWorkspaceDataviewInstances(dataviewInstances)
		if err != nil {
			return nil, err
		}
		workspace.DataviewInstances = &dataviewInstancesObj
	}
	return &workspace, nil
}

func flattenReportWorkspace(workspace api.ReportWorkspace) (map[string]interface{}, error) {
	a := map[string]interface{}{
		"id":          workspace.ID,
		"name":        workspace.Name,
		"description": workspace.Description,
		"app":         workspace.App,
		"category":    workspace.Category,
		"start_at":    workspace.StartAt,
		"end_at":      workspace.EndAt,
	}
	if workspace.Viewport != nil {
		a["viewport"] = []interface{}{flattenWorkspaceViewport(*workspace.Viewport)}
	}
	if workspace.State != nil {
		jsonStr, err := json.Marshal(workspace.State)
		if err != nil {
			return nil, err
		}
		a["state"] = string(jsonStr)
	}
	if workspace.DataviewInstances != nil {
		dataviewInstances, err := flattenWorkspaceDataviewInstances(*workspace.DataviewInstances)
		if err != nil {
			return nil, err
		}
		a["dataview_instances"] = dataviewInstances
	}
	return a, nil
}

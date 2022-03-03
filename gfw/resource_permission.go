package gfw

import (
	"context"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionCreate,
		ReadContext:   resourcePermissionRead,
		UpdateContext: resourcePermissionUpdate,
		DeleteContext: resourcePermissionDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"action": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Required: false,
			},
		},
	}
}

func resourcePermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	action := d.Get("action").([]interface{})[0].(map[string]interface{})
	resource := d.Get("resource").([]interface{})[0].(map[string]interface{})
	description := d.Get("description").(string)

	resourceCreated, err := c.CreatePermission(api.CreatePermission{
		Name:        name,
		Action:      action["id"].(int),
		Resource:    resource["id"].(int),
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(resourceCreated.ID))
	resourcePermissionRead(ctx, d, m)
	return diags
}

func resourcePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	permissionID := d.Id()
	c := m.(*api.GFWClient)
	permission, err := c.GetPermission(permissionID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", permission.Name)
	d.Set("description", permission.Description)
	d.Set("created_at", permission.CreatedAt)

	action := flattenAction(permission.Action)
	if err := d.Set("action", []interface{}{action}); err != nil {
		return diag.FromErr(err)
	}
	resource := flattenResource(permission.Resource)
	if err := d.Set("resource", []interface{}{resource}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePermissionRead(ctx, d, m)
}

func resourcePermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	permissionID := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeletePermission(permissionID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

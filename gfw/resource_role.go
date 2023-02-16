package gfw

import (
	"context"
	"strconv"
	"time"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	actionCreated, err := c.CreateRole(api.CreateRole{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(actionCreated.ID))
	resourceRoleRead(ctx, d, m)
	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	actionId := d.Id()
	c := m.(*api.GFWClient)
	action, err := c.GetRole(actionId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", action.Name)
	d.Set("description", action.Description)
	d.Set("created_at", action.CreatedAt)

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	actionId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteRole(actionId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenRole(action api.Role) interface{} {
	a := make(map[string]interface{})

	a["id"] = action.ID
	a["name"] = action.Name
	a["description"] = action.Description
	a["created_at"] = action.CreatedAt
	return a
}

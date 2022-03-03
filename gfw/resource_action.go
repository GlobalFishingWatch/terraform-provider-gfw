package gfw

import (
	"context"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActionCreate,
		ReadContext:   resourceActionRead,
		UpdateContext: resourceActionUpdate,
		DeleteContext: resourceActionDelete,
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
	}
}

func resourceActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	actionCreated, err := c.CreateAction(api.CreateAction{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(actionCreated.ID))
	resourceActionRead(ctx, d, m)
	return diags
}

func resourceActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	actionId := d.Id()
	c := m.(*api.GFWClient)
	action, err := c.GetAction(actionId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", action.Name)
	d.Set("description", action.Description)
	d.Set("created_at", action.CreatedAt)

	return diags
}

func resourceActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceActionRead(ctx, d, m)
}

func resourceActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	actionId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteAction(actionId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenAction(action api.Action) interface{} {
	a := make(map[string]interface{})

	a["id"] = action.ID
	a["name"] = action.Name
	a["description"] = action.Description
	a["created_at"] = action.CreatedAt
	return a
}

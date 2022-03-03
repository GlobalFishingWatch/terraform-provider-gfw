package gfw

import (
	"context"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceCreate,
		ReadContext:   resourceResourceRead,
		UpdateContext: resourceResourceUpdate,
		DeleteContext: resourceResourceDelete,
		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
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

func resourceResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	rType := d.Get("type").(string)
	rValue := d.Get("value").(string)
	description := d.Get("description").(string)
	resourceCreated, err := c.CreateResource(api.CreateResource{
		Type:        rType,
		Value:       rValue,
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(resourceCreated.ID))
	resourceResourceRead(ctx, d, m)
	return diags
}

func resourceResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceId := d.Id()
	c := m.(*api.GFWClient)
	resource, err := c.GetResource(resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("type", resource.Type)
	d.Set("value", resource.Value)
	d.Set("description", resource.Description)
	d.Set("created_at", resource.CreatedAt)

	return diags
}

func resourceResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceResourceRead(ctx, d, m)
}

func resourceResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	resourceId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteResource(resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenResource(resource api.Resource) interface{} {
	a := make(map[string]interface{})

	a["id"] = resource.ID
	a["type"] = resource.Type
	a["value"] = resource.Value
	a["description"] = resource.Description
	a["created_at"] = resource.CreatedAt
	return a
}

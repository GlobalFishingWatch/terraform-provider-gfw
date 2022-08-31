package gfw

import (
	"context"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	defaultV := d.Get("default").(bool)
	userGroupCreated, err := c.CreateUserGroup(api.CreateUserGroup{
		Name:        name,
		Description: description,
		Default:     defaultV,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(userGroupCreated.ID))
	resourceUserGroupRead(ctx, d, m)
	return diags
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	userGroupId := d.Id()
	c := m.(*api.GFWClient)
	userGroup, err := c.GetUserGroup(userGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", userGroup.Name)
	d.Set("description", userGroup.Description)
	d.Set("default", userGroup.Default)
	d.Set("created_at", userGroup.CreatedAt)

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	defaultV := d.Get("default").(bool)

	userGroupID := d.Id()
	c := m.(*api.GFWClient)
	err := c.UpdateUserGroup(userGroupID, api.CreateUserGroup{
		Name:        name,
		Description: description,
		Default:     defaultV,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceUserGroupRead(ctx, d, m)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	userGroupId := d.Id()

	c := m.(*api.GFWClient)
	_, err := c.DeleteUserGroup(userGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenUserGroup(userGroup api.UserGroup) interface{} {
	a := make(map[string]interface{})

	a["id"] = userGroup.ID
	a["name"] = userGroup.Name
	a["description"] = userGroup.Description
	a["default"] = userGroup.Default
	a["created_at"] = userGroup.CreatedAt
	return a
}

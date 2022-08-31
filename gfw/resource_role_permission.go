package gfw

import (
	"context"
	"strconv"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRolePermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRolePermissionsCreate,
		ReadContext:   resourceRolePermissionsRead,
		UpdateContext: resourceRolePermissionsUpdate,
		DeleteContext: resourceRolePermissionsDelete,
		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Set: schema.HashInt,
			},
		},
	}
}

func resourceRolePermissionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	roleId := d.Get("role").(int)
	permissionIds := utils.ConvertIntSet(d.Get("permissions").(*schema.Set))

	err := c.CreateRolePermissions(api.CreateRolePermissions{
		RoleID:      roleId,
		Permissions: permissionIds,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(roleId))
	resourceRolePermissionsRead(ctx, d, m)
	return diags
}

func resourceRolePermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	roleId := d.Id()
	c := m.(*api.GFWClient)
	role, err := c.GetRole(roleId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("role", role.ID)
	d.Set("permissions", flattenPermissionsToIds(role))

	return diags
}

func resourceRolePermissionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRolePermissionsCreate(ctx, d, m)
}

func resourceRolePermissionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	roleId := d.Id()

	c := m.(*api.GFWClient)
	err := c.DeleteRolePermissions(roleId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenPermissionsToIds(role *api.Role) interface{} {
	var list []int
	for _, p := range role.Permissions {
		list = append(list, p.ID)
	}

	return list
}

package gfw

import (
	"context"
	"strconv"
	"time"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroupRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupRoleCreate,
		ReadContext:   resourceUserGroupRoleRead,
		UpdateContext: resourceUserGroupRoleUpdate,
		DeleteContext: resourceUserGroupRoleDelete,
		Schema: map[string]*schema.Schema{
			"user_group": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Set: schema.HashInt,
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

func resourceUserGroupRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*api.GFWClient)
	var diags diag.Diagnostics

	userGroupID := d.Get("user_group").(int)
	roles := utils.ConvertIntSet(d.Get("roles").(*schema.Set))

	err := c.CreateUserGroupRole(api.CreateUserGroupRole{
		UserGroupID: userGroupID,
		Roles:       roles,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(userGroupID))
	resourceUserGroupRoleRead(ctx, d, m)
	return diags
}

func resourceUserGroupRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	userGroupId := d.Id()
	c := m.(*api.GFWClient)
	userGroup, err := c.GetUserGroup(userGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("user_group", userGroup.ID)
	d.Set("roles", flattenRolesToIds(userGroup))

	return diags
}

func resourceUserGroupRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceUserGroupRoleCreate(ctx, d, m)
}

func resourceUserGroupRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	roleId := d.Id()

	c := m.(*api.GFWClient)
	err := c.DeleteUserGroupRole(roleId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenRolesToIds(userGroup *api.UserGroup) interface{} {
	var list []int
	for _, r := range userGroup.Roles {
		list = append(list, r.ID)
	}

	return list
}

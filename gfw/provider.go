package gfw

import (
	"context"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GFW_TOKEN", nil),
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GFW_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gfw_action":           resourceAction(),
			"gfw_resource":         resourceResource(),
			"gfw_permission":       resourcePermission(),
			"gfw_role":             resourceRole(),
			"gfw_role_permissions": resourceRolePermissions(),
			"gfw_user_group":       resourceUserGroup(),
			"gfw_user_group_role":  resourceUserGroupRole(),
			"gfw_dataset":          resourceDataset(),
			"gfw_dataview":         resourceDataview(),
			"gfw_workspace":        resourceWorkspace(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)

	url := d.Get("url").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := api.NewClient(url, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}

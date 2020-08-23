package main

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

// Provider returns an instance of the
// Confluent Cloud Terraform provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(interface{}, cty.Path) diag.Diagnostics {
					return nil
				},
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateDiagFunc: func(interface{}, cty.Path) diag.Diagnostics {
					return nil
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ccloud_environment": dataSourceEnvironment(),
			"ccloud_cluster":     dataSourceCluster(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ccloud_environment": resourceEnvironment(),
			"ccloud_cluster":     resourceCluster(),
			"ccloud_apikey":      resourceAPIKey(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		session, err := ccloudapi.Login(username, password)
		var diags diag.Diagnostics
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create ccloud client",
				Detail:   "Unable to authenticate user",
			})
			return nil, diags
		}
		return session, nil
	}

	return provider

}

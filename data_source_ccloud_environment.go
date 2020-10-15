package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceEnvironmentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	session := meta.(*ccloudapi.Session)
	name := data.Get("name").(string)
	environments, err := ccloudapi.ListEnvironments(session)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(environments) > 0 {
		for _, environment := range environments {
			if environment.Name == name {
				data.SetId(environment.ID)
				data.Set("name", environment.Name)
				data.Set("organization_id", environment.OrganizationID)
				break
			}
		}
	}
	return diags
}

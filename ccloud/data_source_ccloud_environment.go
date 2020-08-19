package ccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/terraform-provider-ccloud/ccloudapi"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnvironmentRead,
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

func dataSourceEnvironmentRead(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	name := data.Get("name").(string)
	environments, err := ccloudapi.ListEnvironments(session)
	if err != nil {
		return err
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
	return nil
}

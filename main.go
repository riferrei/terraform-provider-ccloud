package main

import (
	"github.com/riferrei/terraform-provider-ccloud/ccloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return ccloud.Provider()
		},
	})
}

package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_ingress": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"network_egress": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"durability": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	session := meta.(*ccloudapi.Session)
	name := data.Get("name").(string)
	environmentID := data.Get("environment_id").(string)
	clusters, err := ccloudapi.ListClusters(environmentID, session)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(clusters) > 0 {
		for _, cluster := range clusters {
			if cluster.Name == name {
				data.SetId(cluster.ID)
				data.Set("name", cluster.Name)
				data.Set("environment_id", cluster.EnvironmentID)
				data.Set("cloud_provider", cluster.CloudProvider)
				data.Set("cloud_region", cluster.CloudRegion)
				data.Set("network_ingress", cluster.NetworkIngress)
				data.Set("network_egress", cluster.NetworkEgress)
				data.Set("storage", cluster.Storage)
				data.Set("durability", cluster.Durability)
				data.Set("organization_id", cluster.OrganizationID)
				data.Set("cluster_endpoint", cluster.ClusterEndpoint)
				data.Set("api_endpoint", cluster.APIEndpoint)
				break
			}
		}
	}
	return diags
}

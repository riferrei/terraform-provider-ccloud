package ccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/terraform-provider-ccloud/ccloudapi"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_ingress": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"network_egress": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"storage": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"durability": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterRead(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	name := data.Get("name").(string)
	environmentID := data.Get("environment_id").(string)
	clusters, err := ccloudapi.ListClusters(environmentID, session)
	if err != nil {
		return err
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
	return nil
}

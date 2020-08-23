package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

const (
	networkIngressLimit = 100
	networkEgressLimit  = 100
	storageLimit        = 5000
)

var (
	cloudProviders = map[string][]string{
		"aws": {
			"ap-southeast-1",
			"eu-central-1",
			"ap-northeast-1",
			"eu-west-3",
			"eu-west-2",
			"us-west-2",
			"eu-west-1",
			"us-east-1",
			"ca-central-1",
			"us-east-2",
			"ap-southeast-2",
			"ap-south-1",
			"us-west-1",
			"sa-east-1"},
		"gcp": {
			"northamerica-northeast1",
			"southamerica-east1",
			"asia-southeast2",
			"us-west4",
			"asia-southeast1",
			"europe-west3",
			"australia-southeast1",
			"us-central1",
			"us-west1",
			"asia-northeast1",
			"us-west2",
			"europe-north1",
			"europe-west4",
			"us-east4",
			"asia-east2",
			"europe-west1",
			"asia-east1",
			"europe-west2",
			"us-east1"},
		"azure": {
			"australiaeast",
			"francecentral",
			"canadacentral",
			"eastus",
			"uksouth",
			"westus2",
			"westeurope",
			"centralus",
			"eastus2",
			"northeurope",
			"southeastasia"},
	}
	selectedCloudProvider string
	durabilityOptions     = []string{"LOW", "HIGH"}
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: clusterCreate,
		Read:   clusterRead,
		Update: clusterUpdate,
		Delete: clusterDelete,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					value, _ := v.(string)
					validProvider := false
					for provider := range cloudProviders {
						if value == provider {
							validProvider = true
							selectedCloudProvider = value
							break
						}
					}
					if !validProvider {
						providers := []string{}
						for provider := range cloudProviders {
							providers = append(providers, provider)
						}
						errors = append(errors, fmt.Errorf("Invalid value for "+
							"cloud provider. Valid values are: "+
							strings.Join(providers, ", ")))
						return warns, errors
					}
					return warns, errors
				},
			},
			"cloud_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					if len(selectedCloudProvider) > 0 {
						value, _ := v.(string)
						validRegion := false
						cloudRegions := cloudProviders[selectedCloudProvider]
						for _, cloudRegion := range cloudRegions {
							if value == cloudRegion {
								validRegion = true
								break
							}
						}
						if !validRegion {
							errors = append(errors, fmt.Errorf("Invalid value for "+
								"cloud region. Valid values are: "+
								strings.Join(cloudRegions, ", ")))
							return warns, errors
						}
					}
					return warns, errors
				},
			},
			"network_ingress": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  networkIngressLimit,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					value, _ := v.(int)
					if value != networkIngressLimit {
						errors = append(errors, fmt.Errorf("Invalid value "+
							"for network ingress. Value needs to be %d",
							networkIngressLimit))
						return warns, errors
					}
					return warns, errors
				},
			},
			"network_egress": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  networkEgressLimit,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					value, _ := v.(int)
					if value != networkEgressLimit {
						errors = append(errors, fmt.Errorf("Invalid value "+
							"for network egress. Value needs to be %d",
							networkEgressLimit))
						return warns, errors
					}
					return warns, errors
				},
			},
			"storage": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  storageLimit,
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warnings []string
					value, _ := v.(int)
					if value != storageLimit {
						errors = append(errors, fmt.Errorf("Invalid value "+
							"for storage. Value needs to be %d",
							storageLimit))
						return warnings, errors
					}
					return warnings, errors
				},
			},
			"durability": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  durabilityOptions[0],
				ValidateFunc: func(v interface{}, k string) (wrs []string, ers []error) {
					var errors []error
					var warns []string
					value, _ := v.(string)
					validDurability := false
					for _, durabilityOption := range durabilityOptions {
						if value == durabilityOption {
							validDurability = true
							break
						}
					}
					if !validDurability {
						errors = append(errors, fmt.Errorf("Invalid value for "+
							"durability. Valid values are: "+
							strings.Join(durabilityOptions, ",")))
						return warns, errors
					}
					return warns, errors
				},
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

func clusterCreate(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	cluster := &ccloudapi.Cluster{
		EnvironmentID:  data.Get("environment_id").(string),
		Name:           data.Get("name").(string),
		CloudProvider:  data.Get("cloud_provider").(string),
		CloudRegion:    data.Get("cloud_region").(string),
		NetworkIngress: data.Get("network_ingress").(int),
		NetworkEgress:  data.Get("network_egress").(int),
		Storage:        data.Get("storage").(int),
		Durability:     data.Get("durability").(string),
	}
	createdCluster, err := ccloudapi.CreateCluster(cluster, session)
	if err != nil {
		return err
	}
	data.SetId(createdCluster.ID)
	data.Set("environment_id", createdCluster.EnvironmentID)
	data.Set("name", createdCluster.Name)
	data.Set("cloud_provider", createdCluster.CloudProvider)
	data.Set("cloud_region", createdCluster.CloudRegion)
	data.Set("network_ingress", createdCluster.NetworkIngress)
	data.Set("network_egress", createdCluster.NetworkEgress)
	data.Set("storage", createdCluster.Storage)
	data.Set("durability", createdCluster.Durability)
	data.Set("organization_id", createdCluster.OrganizationID)
	data.Set("cluster_endpoint", createdCluster.ClusterEndpoint)
	data.Set("api_endpoint", createdCluster.APIEndpoint)
	return nil
}

func clusterRead(data *schema.ResourceData, meta interface{}) error {
	id := data.Id()
	environmentID := data.Get("environment_id").(string)
	session := meta.(*ccloudapi.Session)
	cluster, _ := ccloudapi.ReadCluster(id, environmentID, session)
	if cluster == nil {
		data.SetId("")
		return nil
	}
	data.Set("environment_id", cluster.EnvironmentID)
	data.Set("name", cluster.Name)
	data.Set("cloud_provider", cluster.CloudProvider)
	data.Set("cloud_region", cluster.CloudRegion)
	data.Set("network_ingress", cluster.NetworkIngress)
	data.Set("network_egress", cluster.NetworkEgress)
	data.Set("storage", cluster.Storage)
	data.Set("durability", cluster.Durability)
	data.Set("organization_id", cluster.OrganizationID)
	data.Set("cluster_endpoint", cluster.ClusterEndpoint)
	data.Set("api_endpoint", cluster.APIEndpoint)
	return nil
}

func clusterUpdate(data *schema.ResourceData, meta interface{}) error {
	cluster := &ccloudapi.Cluster{
		ID:             data.Id(),
		EnvironmentID:  data.Get("environment_id").(string),
		Name:           data.Get("name").(string),
		CloudProvider:  data.Get("cloud_provider").(string),
		CloudRegion:    data.Get("cloud_region").(string),
		NetworkIngress: data.Get("network_ingress").(int),
		NetworkEgress:  data.Get("network_egress").(int),
		Storage:        data.Get("storage").(int),
		Durability:     data.Get("durability").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.UpdateCluster(cluster, session)
	return clusterRead(data, meta)
}

func clusterDelete(data *schema.ResourceData, meta interface{}) error {
	cluster := &ccloudapi.Cluster{
		ID:             data.Id(),
		EnvironmentID:  data.Get("environment_id").(string),
		Name:           data.Get("name").(string),
		CloudProvider:  data.Get("cloud_provider").(string),
		CloudRegion:    data.Get("cloud_region").(string),
		NetworkIngress: data.Get("network_ingress").(int),
		NetworkEgress:  data.Get("network_egress").(int),
		Storage:        data.Get("storage").(int),
		Durability:     data.Get("durability").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.DeleteCluster(cluster, session)
	data.SetId("")
	return nil
}

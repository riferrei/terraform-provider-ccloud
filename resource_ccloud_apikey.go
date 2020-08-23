package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/ccloud-sdk-go"
	"strconv"
)

func resourceAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: apiKeyCreate,
		Read:   apiKeyRead,
		Update: apiKeyUpdate,
		Delete: apiKeyDelete,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func apiKeyCreate(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	createdAPIKey, err := ccloudapi.CreateAPIKey(environmentID, clusterID, session)
	if err != nil {
		return err
	}
	data.SetId(strconv.Itoa(createdAPIKey.ID))
	data.Set("key", createdAPIKey.Key)
	data.Set("secret", createdAPIKey.Secret)
	return nil
}

func apiKeyRead(data *schema.ResourceData, meta interface{}) error {
	key := data.Get("key").(string)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	session := meta.(*ccloudapi.Session)
	apiKey, _ := ccloudapi.ReadAPIKey(environmentID, clusterID, key, session)
	if apiKey == nil {
		data.SetId("")
		return nil
	}
	data.Set("key", apiKey.Key)
	return nil
}

func apiKeyUpdate(data *schema.ResourceData, meta interface{}) error {
	return apiKeyRead(data, meta)
}

func apiKeyDelete(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	ccloudapi.DeleteAPIKey(environmentID, clusterID, data.Id(), session)
	data.SetId("")
	return nil
}

package main

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

func resourceAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: apiKeyCreate,
		ReadContext:   apiKeyRead,
		UpdateContext: apiKeyUpdate,
		DeleteContext: apiKeyDelete,
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

func apiKeyCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	session := meta.(*ccloudapi.Session)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	createdAPIKey, err := ccloudapi.CreateAPIKey(environmentID, clusterID, session)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(createdAPIKey.ID))
	data.Set("key", createdAPIKey.Key)
	data.Set("secret", createdAPIKey.Secret)
	return diags
}

func apiKeyRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	key := data.Get("key").(string)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	session := meta.(*ccloudapi.Session)
	apiKey, err := ccloudapi.ReadAPIKey(environmentID, clusterID, key, session)
	if err != nil {
		return diag.FromErr(err)
	}
	if apiKey == nil {
		data.SetId("")
		return diags
	}
	data.Set("key", apiKey.Key)
	return diags
}

func apiKeyUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return apiKeyRead(ctx, data, meta)
}

func apiKeyDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	session := meta.(*ccloudapi.Session)
	environmentID := data.Get("environment_id").(string)
	clusterID := data.Get("cluster_id").(string)
	_, err := ccloudapi.DeleteAPIKey(environmentID, clusterID, data.Id(), session)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	return diags
}

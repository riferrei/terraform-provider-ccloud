package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ccloudapi "github.com/riferrei/ccloud-sdk-go"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: environmentCreate,
		ReadContext:   environmentRead,
		UpdateContext: environmentUpdate,
		DeleteContext: environmentDelete,
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

func environmentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	session := meta.(*ccloudapi.Session)
	environment := &ccloudapi.Environment{
		OrganizationID: session.User.OrganizationID,
		Name:           data.Get("name").(string),
	}
	createdEnvironment, err := ccloudapi.CreateEnvironment(environment, session)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(createdEnvironment.ID)
	data.Set("name", createdEnvironment.Name)
	data.Set("organization_id", createdEnvironment.OrganizationID)
	return diags
}

func environmentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := data.Id()
	session := meta.(*ccloudapi.Session)
	environment, _ := ccloudapi.ReadEnvironment(id, session)
	if environment == nil {
		data.SetId("")
		return diags
	}
	data.Set("name", environment.Name)
	data.Set("organization_id", environment.OrganizationID)
	return diags
}

func environmentUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	environment := &ccloudapi.Environment{
		ID:             data.Id(),
		Name:           data.Get("name").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.UpdateEnvironment(environment, session)
	return environmentRead(ctx, data, meta)
}

func environmentDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	environment := &ccloudapi.Environment{
		ID:             data.Id(),
		Name:           data.Get("name").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.DeleteEnvironment(environment, session)
	data.SetId("")
	return diags
}

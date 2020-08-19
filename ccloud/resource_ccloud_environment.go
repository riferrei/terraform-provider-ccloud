package ccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/terraform-provider-ccloud/ccloudapi"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: environmentCreate,
		Read:   environmentRead,
		Update: environmentUpdate,
		Delete: environmentDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func environmentCreate(data *schema.ResourceData, meta interface{}) error {
	session := meta.(*ccloudapi.Session)
	environment := &ccloudapi.Environment{
		OrganizationID: session.User.OrganizationID,
		Name:           data.Get("name").(string),
	}
	createdEnvironment, err := ccloudapi.CreateEnvironment(environment, session)
	if err != nil {
		return err
	}
	data.SetId(createdEnvironment.ID)
	data.Set("name", createdEnvironment.Name)
	data.Set("organization_id", createdEnvironment.OrganizationID)
	return nil
}

func environmentRead(data *schema.ResourceData, meta interface{}) error {
	id := data.Id()
	session := meta.(*ccloudapi.Session)
	environment, _ := ccloudapi.ReadEnvironment(id, session)
	if environment == nil {
		data.SetId("")
		return nil
	}
	data.Set("name", environment.Name)
	data.Set("organization_id", environment.OrganizationID)
	return nil
}

func environmentUpdate(data *schema.ResourceData, meta interface{}) error {
	environment := &ccloudapi.Environment{
		ID:             data.Id(),
		Name:           data.Get("name").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.UpdateEnvironment(environment, session)
	return environmentRead(data, meta)
}

func environmentDelete(data *schema.ResourceData, meta interface{}) error {
	environment := &ccloudapi.Environment{
		ID:             data.Id(),
		Name:           data.Get("name").(string),
		OrganizationID: data.Get("organization_id").(int),
	}
	session := meta.(*ccloudapi.Session)
	ccloudapi.DeleteEnvironment(environment, session)
	data.SetId("")
	return nil
}

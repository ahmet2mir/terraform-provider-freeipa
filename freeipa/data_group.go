package freeipa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			groupSchemaGid: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			groupSchemaGidNumber: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			groupSchemaDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, m interface{}) error {

	groupIDData, groupIDOk := d.GetOk(groupSchemaGid)
	log.Printf("[TRACE] dataSourceGroupRead - %s", groupIDData)

	if !groupIDOk {
		log.Printf("[ERROR] dataSourceGroupRead - one and only one of %s must be set", groupSchemaGid)
		return nil
	}

	c := m.(*Connection)

	gid, err := c.ldapClient.GetGroupForGroupname(groupIDData.(string))

	if err != nil {
		log.Printf("[ERROR] Error reading group %s - %s", gid, err)
		return err
	}

	rec, err := c.client.GetGroup(*gid)

	if err != nil {
		log.Printf("[ERROR] Error getting group %s - %s", gid, err)
		return err
	}

	d.SetId(fmt.Sprintf("%d", rec.Gid))

	err = d.Set(groupSchemaGid, rec.Gid)
	if err != nil {
		return err
	}

	err = d.Set(groupSchemaDescription, rec.Description)
	if err != nil {
		return err
	}

	err = d.Set(groupSchemaGidNumber, rec.GidNumber)
	if err != nil {
		return err
	}
	return nil
}

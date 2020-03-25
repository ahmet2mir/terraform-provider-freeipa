package freeipa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			userSchemaUid: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			userSchemaEmail: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaFirstName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			userSchemaLastName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			userSchemaUidNumber: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaGidNumber: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			userSchemaGroups: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {

	userIDData, userIDOk := d.GetOk(userSchemaUid)
	log.Printf("[TRACE] dataSourceUserRead - %s", userIDData)

	if !userIDOk {
		log.Printf("[ERROR] dataSourceUserRead - one and only one of %s must be set", userSchemaUid)
		return nil
	}

	c := m.(*Connection)

	uid, err := c.ldapClient.GetUserForUsername(userIDData.(string))

	if err != nil {
		log.Printf("[ERROR] dataSourceUserRead - Error reading user %s - %s", uid, err)
		return err
	}

	rec, err := c.client.GetUser(*uid)

	if err != nil {
		log.Printf("[ERROR] dataSourceUserRead - Error getting user %s - %s", uid, err)
		return err
	}

	d.SetId(fmt.Sprintf("%d", rec.Uid))

	err = d.Set(userSchemaUid, rec.Uid)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaUidNumber, rec.UidNumber)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaLastName, rec.Last)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaFirstName, rec.First)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaEmail, rec.Email)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaGidNumber, rec.GidNumber)
	if err != nil {
		return err
	}

	err = d.Set(userSchemaGroups, rec.Groups)
	if err != nil {
		return err
	}

	return nil
}

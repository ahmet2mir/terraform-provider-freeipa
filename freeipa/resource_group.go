package freeipa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

const (
	groupSchemaGid          = "gid"
	groupSchemaGidNumber    = "gid_number"
	groupSchemaDescription  = "description"
	groupSchemaGroups       = "groups"
	groupSchemaGroupMembers = "group_members"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			groupSchemaGroups: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			groupSchemaGroupMembers: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceGroupCreate - %s", d.Id())

	c := m.(*Connection)

	gid := d.Get(groupSchemaGid).(string)
	description := d.Get(groupSchemaDescription).(string)

	log.Printf("[TRACE] creating group with name - %s, description - %s",
		gid, description)

	options := map[string]interface{}{}

	gidNumber, ok := d.GetOk(groupSchemaGidNumber)
	if ok {
		options["gidnumber"] = gidNumber.(string)
	}

	groupRec, err := c.client.CreateGroup(gid, description, options)

	if err != nil {
		log.Printf("[ERROR] Error creating group %s - %s", gid, err)
		return err
	}

	d.SetId(string(groupRec.IpaUniqueId))

	groupsInterface, ok := d.GetOk(groupSchemaGroups)
	if ok {
		groupsRaw := groupsInterface.(*schema.Set)
		if groupsRaw.Len() > 0 {
			groups := make([]string, groupsRaw.Len())
			for i, d := range groupsRaw.List() {
				groups[i] = d.(string)
			}

			err = c.client.GroupSyncGroups(gid, groups, false)
			if err != nil {
				return err
			}
		}
	}

	groupMembersInterface, ok := d.GetOk(groupSchemaGroupMembers)
	if ok {
		groupsRaw := groupMembersInterface.(*schema.Set)
		if groupsRaw.Len() > 0 {
			groups := make([]string, groupsRaw.Len())
			for i, d := range groupsRaw.List() {
				groups[i] = d.(string)
			}

			err = c.client.GroupSyncGroups(gid, groups, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceGroupRead - %s", d.Id())

	c := m.(*Connection)

	gid, err := c.ldapClient.GetGroupForUUID(d.Id())

	if err != nil {
		log.Printf("[ERROR] Error reading group %s - %s", gid, err)
		return err
	}

	rec, err := c.client.GetGroup(*gid)

	if err != nil {
		log.Printf("[ERROR] Error getting group %s - %s", gid, err)
		return err
	}

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

	err = d.Set(groupSchemaGroups, rec.Groups)
	if err != nil {
		return err
	}

	err = d.Set(groupSchemaGroupMembers, rec.GroupMembers)
	if err != nil {
		return err
	}
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceGroupUpdate - %s", d.Id())

	c := m.(*Connection)

	gid, err := c.ldapClient.GetGroupForUUID(d.Id())

	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange(groupSchemaGid) {
		_, newValue := d.GetChange(groupSchemaGid)
		val := newValue.(string)
		c.client.GroupUpdateGid(*gid, val)
		d.SetPartial(groupSchemaGid)
		gid = &val
	}

	if d.HasChange(groupSchemaGidNumber) {
		_, newValue := d.GetChange(groupSchemaGidNumber)
		c.client.GroupUpdateGidNumber(*gid, newValue.(string))
		d.SetPartial(groupSchemaGidNumber)
	}

	if d.HasChange(groupSchemaDescription) {
		_, newValue := d.GetChange(groupSchemaDescription)
		c.client.GroupUpdateDescription(*gid, newValue.(string))
		d.SetPartial(groupSchemaDescription)
	}

	if d.HasChange(groupSchemaGroups) {
		_, newValueInterface := d.GetChange(groupSchemaGroups)

		groupsRaw := newValueInterface.(*schema.Set)
		newValue := make([]string, groupsRaw.Len())
		for i, d := range groupsRaw.List() {
			newValue[i] = d.(string)
		}

		c.client.GroupSyncGroups(*gid, newValue, false)
		d.SetPartial(groupSchemaGroups)
	}

	if d.HasChange(groupSchemaGroupMembers) {
		_, newValueInterface := d.GetChange(groupSchemaGroupMembers)

		groupsRaw := newValueInterface.(*schema.Set)
		newValue := make([]string, groupsRaw.Len())
		for i, d := range groupsRaw.List() {
			newValue[i] = d.(string)
		}

		c.client.GroupSyncGroups(*gid, newValue, true)
		d.SetPartial(groupSchemaGroupMembers)
	}

	d.Partial(false)

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceGroupDelete - %s", d.Id())

	c := m.(*Connection)

	gid, err := c.ldapClient.GetGroupForUUID(d.Id())

	if err != nil {
		return err
	}

	return c.client.DeleteGroup(*gid)
}

func resourceGroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Printf("[TRACE] resourceGroupExists - %s", d.Id())

	c := m.(*Connection)

	id := d.Id()

	exists, err := c.ldapClient.GroupExistsForUUID(id)

	if err != nil {
		return false, err
	}

	return exists, nil
}

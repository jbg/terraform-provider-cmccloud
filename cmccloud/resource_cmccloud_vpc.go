package cmccloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudVPCCreate,
		Read:   resourceCMCCloudVPCRead,
		Update: resourceCMCCloudVPCUpdate,
		Delete: resourceCMCCloudVPCDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudVPCImport,
		},
		SchemaVersion: 1,
		Schema:        vpcSchema(),
	}
}

func resourceCMCCloudVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, taskStatus, err := client.VPC.Create(d.Get("name").(string), d.Get("description").(string), d.Get("region").(string), d.Get("cidr").(string))
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}
	d.SetId(taskStatus.ResultID)

	return resourceCMCCloudVPCRead(d, meta)
}

func resourceCMCCloudVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc, err := client.VPC.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving VPC %s: %v", d.Id(), err)
	}

	_ = d.Set("name", vpc.Name)
	_ = d.Set("description", vpc.Description)
	_ = d.Set("state", vpc.State)
	_ = d.Set("region", strings.ToLower(vpc.RegionName))
	_ = d.Set("cidr", vpc.Cidr)
	return nil
}

func resourceCMCCloudVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		err := client.VPC.Update(id, d.Get("name").(string), d.Get("description").(string))
		if err != nil {
			return fmt.Errorf("Error when rename VPC [%s]: %v", id, err)
		}
	}
	return resourceCMCCloudVPCRead(d, meta)
}

func resourceCMCCloudVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.VPC.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud vpc: %v", err)
	}
	return nil
}

func resourceCMCCloudVPCImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCMCCloudVPCRead(d, meta)
	return []*schema.ResourceData{d}, err
}

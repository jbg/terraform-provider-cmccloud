package cmccloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudFloatingIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudFloatingIPCreate,
		Read:   resourceCMCCloudFloatingIPRead,
		// Update: resourceCMCCloudFloatingIPUpdate,
		Delete: resourceCMCCloudFloatingIPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudFloatingIPImport,
		},
		SchemaVersion: 1,
		Schema:        floatingIPSchema(),
	}
}

func resourceCMCCloudFloatingIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, taskStatus, err := client.FloatingIP.Create(d.Get("vpc_id").(string))
	if err != nil {
		return fmt.Errorf("Error creating FloatingIP: %s", err)
	}
	d.SetId(taskStatus.ResultID)

	return resourceCMCCloudFloatingIPRead(d, meta)
}

func resourceCMCCloudFloatingIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	FloatingIP, err := client.FloatingIP.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving FloatingIP %s: %v", d.Id(), err)
	}

	_ = d.Set("vpc_id", FloatingIP.VPCID)
	_ = d.Set("ip_address", FloatingIP.IPAddress)
	_ = d.Set("is_source_nat", FloatingIP.IsSourceNat)
	_ = d.Set("is_static_nat", FloatingIP.IsStaticNat)
	_ = d.Set("state", FloatingIP.State)

	return nil
}

func resourceCMCCloudFloatingIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.FloatingIP.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud FloatingIP: %v", err)
	}
	return nil
}

func resourceCMCCloudFloatingIPImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCMCCloudFloatingIPRead(d, meta)
	return []*schema.ResourceData{d}, err
}

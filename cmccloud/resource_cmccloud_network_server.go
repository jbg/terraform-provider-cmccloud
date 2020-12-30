package cmccloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudNetworkServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudNetworkServerCreate,
		Read:   resourceCMCCloudNetworkServerRead,
		Delete: resourceCMCCloudNetworkServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudNetworkServerImport,
		},
		SchemaVersion: 1,
		Schema:        networkServerSchema(),
	}
}

func resourceCMCCloudNetworkServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	networkID := d.Get("network_id").(string)
	serverID := d.Get("server_id").(string)

	_, err := client.Server.AddNic(serverID, networkID)
	if err != nil {
		return fmt.Errorf("Error creating Network: %s", err)
	}
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-%s", serverID, networkID)))
	return resourceCMCCloudNetworkRead(d, meta)
}

func resourceCMCCloudNetworkServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	serverID := d.Get("server_id").(string)
	networkID := d.Get("network_id").(string)

	server, err := client.Server.Get(serverID)
	if err != nil {
		d.SetId("")
		return nil
	}

	if server.Nics != nil {
		found := false
		for _, nic := range server.Nics {
			if nic.NetworkID == networkID {
				found = true
			}
		}
		if !found {
			log.Printf("[DEBUG] Not found this network %s on this server %s", networkID, serverID)
			d.SetId("")
		}

	} else {
		log.Printf("[DEBUG] Not found any nic on this server %s", serverID)
		d.SetId("")
	}

	return nil
}

func resourceCMCCloudNetworkServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Network.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud Network: %v", err)
	}
	return nil
}

func resourceCMCCloudNetworkServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceCMCCloudNetworkRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

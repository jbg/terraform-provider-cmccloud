package cmccloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudNetworkCreate,
		Read:   resourceCMCCloudNetworkRead,
		Update: resourceCMCCloudNetworkUpdate,
		Delete: resourceCMCCloudNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudNetworkImport,
		},
		SchemaVersion: 1,
		Schema:        networkSchema(),
	}
}

func resourceCMCCloudNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	result, err := client.Network.CreateVPCNetwork(d.Get("vpc_id").(string),
		d.Get("name").(string),
		d.Get("description").(string),
		d.Get("gateway").(string),
		d.Get("netmask").(string),
		d.Get("firewall_id").(string))
	if err != nil {
		return fmt.Errorf("Error creating Network: %s", err)
	}
	d.SetId(result.ResultID)
	return resourceCMCCloudNetworkRead(d, meta)
}

func resourceCMCCloudNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	network, err := client.Network.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Network %s: %v", d.Id(), err)
	}

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("gateway", network.Gateway)
	d.Set("netmask", network.Netmask)
	d.Set("cidr", network.Cidr)
	d.Set("state", network.State)
	d.Set("type", network.Type)
	d.Set("firewall_id", fixFirewallID(network.FirewallID))
	d.Set("vpc_id", network.VPCID)
	d.Set("server_ids", stringArrayToSet(network.ServerIDs))

	return nil
}

func fixFirewallID(firewallID string) string {
	if firewallID == "96552912-1e2d-11e6-9f67-062a7b070809" {
		return "allow"
	} else if firewallID == "9654d282-1e2d-11e6-9f67-062a7b070809" {
		return "deny"
	} else {
		return firewallID
	}
}

func resourceCMCCloudNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		err := client.Network.Update(id, d.Get("name").(string), d.Get("description").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Network [%s]: %v", id, err)
		}
	}

	if d.HasChange("firewall_id") {
		_, err := client.Network.ChangeFirewall(id, d.Get("firewall_id").(string))
		if err != nil {
			return fmt.Errorf("Error when change firewall of Network [%s]: to [%s] %v", id, d.Get("firewall_id"), err)
		}
	}

	if d.HasChange("server_ids") {
		oldIDs, newIDs := d.GetChange("server_ids")
		newSet := func(ids []interface{}) map[string]struct{} {
			out := make(map[string]struct{}, len(ids))
			for _, id := range ids {
				out[id.(string)] = struct{}{}
			}
			return out
		}
		// leftDiff returns all elements in Left that are not in Right
		leftDiff := func(left, right map[string]struct{}) map[string]struct{} {
			out := make(map[string]struct{})
			for l := range left {
				if _, ok := right[l]; !ok {
					out[l] = struct{}{}
				}
			}
			return out
		}
		oldIDSet := newSet(oldIDs.(*schema.Set).List())
		newIDSet := newSet(newIDs.(*schema.Set).List())
		for serverID := range leftDiff(newIDSet, oldIDSet) {
			_, err := client.Server.AddNic(serverID, id)
			if err != nil {
				return fmt.Errorf("Error when create attach network %s to server (%s): %s", d.Id(), serverID, err)
			}
		}
		for serverID := range leftDiff(oldIDSet, newIDSet) {
			server, err := client.Server.Get(serverID)
			if err != nil {
				return fmt.Errorf("Error when get server (%s) info: %s", serverID, err)
			}
			nics := server.Nics
			for _, nic := range nics {
				if nic.NetworkID == id {
					_, err := client.Server.RemoveNic(serverID, nic.ID)
					if err != nil {
						return fmt.Errorf("Error when create remove nic %s (belong to network %s) from server (%s): %s", nic.ID, d.Id(), serverID, err)
					}
					break
				}
			}
		}
	}
	return resourceCMCCloudNetworkRead(d, meta)
}

func resourceCMCCloudNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Network.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud Network: %v", err)
	}
	return nil
}

func resourceCMCCloudNetworkImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceCMCCloudNetworkRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

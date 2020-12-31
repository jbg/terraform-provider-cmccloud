package cmccloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cmc-cloud/gocmcapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudFirewallVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudFirewallVPCCreate,
		Read:   resourceCMCCloudFirewallVPCRead,
		Update: resourceCMCCloudFirewallVPCUpdate,
		Delete: resourceCMCCloudFirewallVPCDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudFirewallVPCImport,
		},
		SchemaVersion: 1,
		Schema:        firewallVPCSchema(),
	}
}

func resourceCMCCloudFirewallVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	taskStatus, err := client.FirewallVPC.Create(d.Get("vpc_id").(string), d.Get("name").(string), d.Get("description").(string))
	if err != nil {
		return fmt.Errorf("Error creating FirewallVPC: %s", err)
	}
	id := taskStatus.ResultID
	d.SetId(id)
	return resourceCMCCloudFirewallVPCUpdate(d, meta)
}

func resourceCMCCloudFirewallVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	firewallVPC, err := client.FirewallVPC.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving FirewallVPC %s: %v", d.Id(), err)
	}

	d.Set("name", firewallVPC.Name)
	d.Set("description", firewallVPC.Description)
	d.Set("vpc_id", firewallVPC.VPCID)
	d.Set("inbound_rule", convertFirewallRule(firewallVPC.InboundRules))
	d.Set("outbound_rule", convertFirewallRule(firewallVPC.OutboundRules))
	return nil
}

func resourceCMCCloudFirewallVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	inboundJSON := "[]"
	outboundJSON := "[]"

	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		err := client.FirewallVPC.Update(id, d.Get("name").(string), d.Get("description").(string))
		if err != nil {
			return fmt.Errorf("Error when rename FirewallVPC [%s]: %v", id, err)
		}
	}
	if d.HasChange("inbound_rule") || d.HasChange("outbound_rule") {
		if v, ok := d.GetOk("inbound_rule"); ok {
			res, err := json.Marshal(flatternFirewallRules(v.([]interface{})))
			if err != nil {
				return fmt.Errorf("Error when marshal inbound rules: %v, %+v", err, v)
			}
			inboundJSON = string(res)
		}
		if v, ok := d.GetOk("outbound_rule"); ok {
			res, err := json.Marshal(flatternFirewallRules(v.([]interface{})))
			if err != nil {
				return fmt.Errorf("Error when marshal outbound rules: %v, %+v", err, v)
			}
			outboundJSON = string(res)
		}

		_, err := client.FirewallVPC.SaveRules(id, inboundJSON, outboundJSON)
		if err != nil {
			return fmt.Errorf("Error when update rules: %v", err)
		}
	}
	return resourceCMCCloudFirewallVPCRead(d, meta)
}

func resourceCMCCloudFirewallVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.FirewallVPC.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud FirewallVPC: %v", err)
	}
	return nil
}

func resourceCMCCloudFirewallVPCImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCMCCloudFirewallVPCRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func convertFirewallRule(rules []gocmcapi.FirewallVPCRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		result[i] = map[string]interface{}{
			"id":         rule.ID,
			"cidrs":      rule.Cidrs,
			"action":     strings.ToLower(rule.Action),
			"port_range": rule.PortRange,
			"protocol":   rule.Protocol,
		}
	}
	return result
}

func flatternFirewallRules(rules []interface{}) []gocmcapi.FirewallVPCRule {
	fwrules := []gocmcapi.FirewallVPCRule{}
	for _, rawRule := range rules {
		r := rawRule.(map[string]interface{})
		rule := gocmcapi.FirewallVPCRule{
			ID:        r["id"].(string),
			Protocol:  r["protocol"].(string),
			Action:    r["action"].(string),
			Cidrs:     setToStringArray(r["cidrs"].(*schema.Set)),
			PortRange: r["port_range"].(string),
		}
		fwrules = append(fwrules, rule)
	}
	return fwrules
}

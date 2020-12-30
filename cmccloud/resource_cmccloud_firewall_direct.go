package cmccloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cmccloud/gocmcapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCMCCloudFirewallDirect() *schema.Resource {
	return &schema.Resource{
		Create: resourceCMCCloudFirewallDirectCreate,
		Read:   resourceCMCCloudFirewallDirectRead,
		Update: resourceCMCCloudFirewallDirectUpdate,
		Delete: resourceCMCCloudFirewallDirectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCMCCloudFirewallDirectImport,
		},
		SchemaVersion: 1,
		Schema:        firewallDirectSchema(),
	}
}

func resourceCMCCloudFirewallDirectCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("server_id").(string), d.Get("ip_address").(string)))
	return resourceCMCCloudFirewallDirectUpdate(d, meta)
}

func resourceCMCCloudFirewallDirectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if len(strings.Split(id, ":")) == 2 {
		serverID := strings.Split(id, ":")[0]
		ipAddress := strings.Split(id, ":")[1]

		FirewallDirect, err := client.FirewallDirect.Get(serverID, ipAddress)
		if err != nil {
			return fmt.Errorf("Error receiving FirewallDirect %s: %v", d.Id(), err)
		}
		d.Set("server_id", FirewallDirect.ServerID)
		d.Set("ip_address", FirewallDirect.IPAddress)
		d.Set("inbound_rule", convertFirewallDirectRule(FirewallDirect.InboundRules))
		d.Set("outbound_rule", convertFirewallDirectRule(FirewallDirect.OutboundRules))
	} else {
		d.SetId("")
		return fmt.Errorf("Error get FirewallDirect: invalid id %v", id)
	}

	return nil
}

func resourceCMCCloudFirewallDirectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	inboundJSON := "[]"
	outboundJSON := "[]"

	if d.HasChange("inbound_rule") || d.HasChange("outbound_rule") {
		if v, ok := d.GetOk("inbound_rule"); ok && v != nil {
			res, err := json.Marshal(flatternFirewallDirectRules(v.([]interface{})))
			if err != nil {
				return fmt.Errorf("Error when marshal inbound rules: %v, %+v", err, v)
			}
			inboundJSON = string(res)
		}
		if v, ok := d.GetOk("outbound_rule"); ok && v != nil {
			res, err := json.Marshal(flatternFirewallDirectRules(v.([]interface{})))
			if err != nil {
				return fmt.Errorf("Error when marshal outbound rules: %v, %+v", err, v)
			}
			outboundJSON = string(res)
		}

		_, err := client.FirewallDirect.SaveRules(d.Get("server_id").(string), d.Get("ip_address").(string), inboundJSON, outboundJSON)
		if err != nil {
			return fmt.Errorf("Error when update rules: %v", err)
		}
	}
	return resourceCMCCloudFirewallDirectRead(d, meta)
}

func resourceCMCCloudFirewallDirectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if len(strings.Split(id, ":")) == 2 {
		serverID := strings.Split(id, ":")[0]
		ipAddress := strings.Split(id, ":")[1]
		_, err := client.FirewallDirect.Delete(serverID, ipAddress)
		if err != nil {
			return fmt.Errorf("Error delete cloud FirewallDirect: %v", err)
		}
	} else {
		return fmt.Errorf("Error delete cloud FirewallDirect: invalid id %v", id)
	}
	return nil
}

func resourceCMCCloudFirewallDirectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceCMCCloudFirewallDirectRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

func convertFirewallDirectRule(rules []gocmcapi.FirewallDirectRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		result[i] = map[string]interface{}{
			"id":         rule.ID,
			"name":       rule.Name,
			"protocol":   rule.Protocol,
			"action":     strings.ToLower(rule.Action),
			"port_range": rule.PortRange,
			"port_type":  rule.PortType,
		}
		if rule.Src != "" {
			result[i]["src"] = rule.Src
		}

		if rule.Dst != "" {
			result[i]["dst"] = rule.Dst
		}
	}
	return result
}

func flatternFirewallDirectRules(rules []interface{}) []gocmcapi.FirewallDirectRule {
	fwrules := []gocmcapi.FirewallDirectRule{}
	for _, rawRule := range rules {
		r := rawRule.(map[string]interface{})
		rule := gocmcapi.FirewallDirectRule{
			ID:        r["id"].(int),
			Name:      r["name"].(string),
			Protocol:  r["protocol"].(string),
			Action:    r["action"].(string),
			PortRange: r["port_range"].(string),
			PortType:  r["port_type"].(string),
		}
		if src, ok := r["src"]; ok {
			rule.Src = src.(string)
		}

		if dst, ok := r["dst"]; ok {
			rule.Dst = dst.(string)
		}
		fwrules = append(fwrules, rule)
	}
	return fwrules
}

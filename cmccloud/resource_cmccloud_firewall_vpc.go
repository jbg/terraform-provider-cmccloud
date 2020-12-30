package cmccloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cmccloud/gocmcapi"
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

	// Get inbound_rules
	/*number := 0
	if v, ok := d.GetOk("inbound_rule"); ok {
		number++
		rules := v.([]interface{})
		for _, rawRule := range rules {
			rule := rawRule.(map[string]interface{})
			action := "allow"
			if v, ok := rule["allow"]; ok {
				action = v.(string)
			}
			portRange := ""
			if v, ok := rule["port_range"]; ok {
				portRange = v.(string)
			}
			cidrs := strings.Join(setToStringArray(rule["cidrs"].(*schema.Set)), ",")
			client.FirewallVPC.CreateRule(id, number, cidrs, action, rule["protocol"].(string), "ingress", portRange)
		}
	}

	maxNumber := 0
	if v, ok := d.GetOk("inbound_rule"); ok {
		err := saveRules(client, id, v, "ingress", maxNumber)
		if err != "" {
			return fmt.Errorf(err)
		}
		maxNumber += len(v.([]interface{}))
	}

	if v, ok := d.GetOk("outbound_rule"); ok {
		err := saveRules(client, id, v, "ingress", maxNumber)
		if err != "" {
			return fmt.Errorf(err)
		}
		maxNumber += len(v.([]interface{}))
	}*/
	return resourceCMCCloudFirewallVPCUpdate(d, meta)

	//return resourceCMCCloudFirewallVPCRead(d, meta)
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
	/*
		errors := validateRules(d, meta)
		if errors != nil && len(errors) > 0 {
			return fmt.Errorf("Error when validate rules: %v", errors)
		}

		maxNumber, err := getCurrentMaxNumber(d, meta)
		if err != nil {
			return fmt.Errorf("Error when getting list of rules of this firewall: %v", err)
		}
		err = client.FirewallVPC.DeleteAllRules(id)
		if err != nil {
			return fmt.Errorf("Error when delete all previous rules: %v", err)
		}
		if d.HasChange("inbound_rule") {
			if v, ok := d.GetOk("inbound_rule"); ok {
				err := saveRules(client, id, v, "ingress", maxNumber)
				if err != "" {
					return fmt.Errorf(err)
				}
				maxNumber += len(v.([]interface{}))
			}
		}
		if d.HasChange("outbound_rule") {
			if v, ok := d.GetOk("outbound_rule"); ok {
				err := saveRules(client, id, v, "egress", maxNumber)
				if err != "" {
					return fmt.Errorf(err)
				}
			}
		} */
	return resourceCMCCloudFirewallVPCRead(d, meta)
}

func validateRules(d *schema.ResourceData, meta interface{}) []string {
	client := meta.(*CombinedConfig).goCMCClient()
	inboundJSON := "[]"
	outboundJSON := "[]"
	if v, ok := d.GetOk("inbound_rule"); ok {
		res, err := json.Marshal(flatternFirewallRules(v.([]interface{})))
		if err != nil {
			return []string{fmt.Sprintf("Error when validate rules: %v, %+v", err, v)}
		}
		inboundJSON = string(res)
	}
	if v, ok := d.GetOk("outbound_rule"); ok {
		res, err := json.Marshal(flatternFirewallRules(v.([]interface{})))
		if err != nil {
			return []string{fmt.Sprintf("Error when validate rules: %v, %+v", err, v)}
		}
		outboundJSON = string(res)
	}
	errors, err := client.FirewallVPC.ValidateRules(inboundJSON, outboundJSON)
	if err != nil {
		return []string{fmt.Sprintf("Error when validate rules: %v", err)}
	}
	if errors != nil && len(errors) > 0 {
		return errors
	}
	return nil
}

func saveRules(client *gocmcapi.Client, firewallID string, v interface{}, ruleType string, startNumber int) string {
	rules := v.([]interface{})
	number := startNumber
	for _, rawRule := range rules {
		number++
		rule := rawRule.(map[string]interface{})
		action := "allow"
		if v, ok := rule["allow"]; ok {
			action = v.(string)
		}
		portRange := ""
		if v, ok := rule["port_range"]; ok {
			portRange = v.(string)
		}
		cidrs := strings.Join(setToStringArray(rule["cidrs"].(*schema.Set)), ",")
		if ruleID, ok := rule["id"]; ok && ruleID != "" {
			_, err := client.FirewallVPC.UpdateRule(ruleID.(string), number, cidrs, action, rule["protocol"].(string), ruleType, portRange)
			if err != nil {
				return fmt.Sprintf("Error when update firewall with id = %s: %v", ruleID, err)
			}
		} else {
			_, err := client.FirewallVPC.CreateRule(firewallID, number, cidrs, action, rule["protocol"].(string), ruleType, portRange)
			if err != nil {
				return fmt.Sprintf("Error when creating firewall with info = %v: %v", rule, err)
			}
		}
	}
	return ""
}
func getCurrentMaxNumber(d *schema.ResourceData, meta interface{}) (int, error) {
	client := meta.(*CombinedConfig).goCMCClient()
	rules, err := client.FirewallVPC.GetRules(d.Id())

	if err != nil {
		return 0, err
	}
	maxNumber := 0
	for _, rawRule := range rules {
		rule := rawRule.(map[string]interface{})
		number := int(rule["number"].(float64))
		if maxNumber < number {
			maxNumber = number
		}
	}
	return maxNumber, nil
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
	resourceCMCCloudFirewallVPCRead(d, meta)
	return []*schema.ResourceData{d}, nil
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

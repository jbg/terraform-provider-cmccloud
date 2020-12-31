package cmccloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ruleSchema(cidrName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"protocol": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"tcp",
				"udp",
				"icmp",
			}, false),
		},
		"action": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "allow",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.ToLower(old) == strings.ToLower(new)
			},
			StateFunc: func(val interface{}) string {
				if val == nil || val.(string) == "" {
					return "allow"
				}
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{
				"allow",
				"deny",
			}, false),
		},
		"port_range": {
			Type:     schema.TypeString,
			Optional: true,
		},
		cidrName: {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateIPCidrRange,
			},
			Optional: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func directRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"src": {
			Type:         schema.TypeString,
			ValidateFunc: validateIPCidrRange,
			Optional:     true,
		},
		"dst": {
			Type:         schema.TypeString,
			ValidateFunc: validateIPAddress,
			Optional:     true,
		},
		"action": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "allow",
			ValidateFunc: validation.StringInSlice([]string{
				"allow", "block",
			}, false),
		},
		"port_range": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "local",
			ValidateFunc: validation.StringInSlice([]string{
				"local", "remote",
			}, false),
		},
		"protocol": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"tcp",
				"udp",
				"icmp",
				"ip",
			}, false),
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}
func firewallDirectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"inbound_rule": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: directRuleSchema(),
			},
		},
		"outbound_rule": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: directRuleSchema(),
			},
		},
	}
}

func firewallVPCSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"vpc_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"inbound_rule": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Optional: true,
			//Computed:   true,
			//ConfigMode: schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: ruleSchema("cidrs"),
			},
		},
		"outbound_rule": {
			Type:     schema.TypeList,
			Optional: true,
			//Computed:   true,
			//ConfigMode: schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: ruleSchema("cidrs"),
			},
		},
	}
}

func networkServerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"network_id": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},

		"server_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
	}
}

func networkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"gateway": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateIPAddress,
		},
		"netmask": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateNetmask,
		},
		"vpc_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"firewall_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateFirewallID,
		},
		"server_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			// DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// 	 if jsonold == jsonnew {
			// 		return true
			// 	}
			// 	return false
			// },
		},
		"cidr": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func floatingIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_source_nat": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"is_static_nat": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func vpcSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{"hn", "hcm"}, true),
		},
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
			ValidateFunc: validateIPCidrRange,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func volumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateRegexp(`(?m)^\d{1,10}0$`), // size must be end with 0
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{"ssd", "hdd"}, true),
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{"hn", "hcm", ""}, true),
		},
		"server_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
		},
	}
}

func snapshotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"volume_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
		},
		"server_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"volume_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func serverSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"cpu": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"root": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"gpu": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
		},
		"image_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{"image", "snapshot", "backup"}, true),
		},
		"image_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.ToLower(val.(string))
			},
			ValidateFunc: validation.StringInSlice([]string{"hn", "hcm"}, true),
		},
		"enable_private_network": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"auto_backup": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"backup_schedule": {
			Type:     schema.TypeString,
			Optional: true,
			StateFunc: func(val interface{}) string {
				return strings.ToUpper(val.(string))
			},
			ValidateFunc: validateRegexp(`(?i)^(DAILY|WEEKLY|MONTHLY)-[0-1][0-9]:[0-5][0-9](:\d{1,2})?$`),
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bits": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"main_ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

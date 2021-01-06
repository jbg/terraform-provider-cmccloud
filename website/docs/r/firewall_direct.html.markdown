---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_firewall_direct"
sidebar_current: "docs-cmccloud-resource-firewall-direct"
description: |-
  Provides a CMC Cloud Firewall Direct resource. This can be used to create, modify, and delete Firewall Direct. Firewall Direct is a firewall for a specific IP address of a server
---

# cmccloud\_firewall_direct

Provides a CMC Cloud Firewall Direct resource. This can be used to create,
modify, and delete Firewall Direct. This firewall will be used for networks in Direct

## Example Create a Firewall Direct

```hcl
# Create a new Firewall for Direct
resource "cmccloud_firewall_direct" "terra_firewall_2" {
    server_id = "040e6918-409e-4346-8827-cd153e113f91" 
    ip_address = "203.171.21.11"
    inbound_rule {
        name = "allow ssh"
        protocol         = "tcp"
        port_range       = "22"
        port_type = "local"
        src  = "0.0.0.0/0"
        action = "allow"
    }
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The ID of Server
* `ip_address` - (Required) The ip address of nic that will be used for apply firewall
* `inbound_rule` - (Optional) Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below
* `outbound_rule` - (Optional) Can be specified multiple times for each egress rule. Each egress block supports fields documented below.

The `inbound_rule` block support the following: 
* `name` - (Required) CIDR Block: IPv4 or IPv6 CIDR
* `src` - (Required) Source address cidr
* `action` - (Optional) Port type, Available: allow, block
* `port_range` - (Optional) Port range. Example: `80` or `8000-9000`
* `port_type` - (Optional) Port type, Available: local, remote
* `protocol` - (Required) Layer 4 protocol.  Available: tcp, udp, icmp, ip

The `outbound_rule` block support the following: 
* `name` - (Required) CIDR Block: IPv4 or IPv6 CIDR 
* `dst` - (Required) Destination address cidr
* `action` - (Optional) Port type, Available: allow, block
* `port_range` - (Optional) Port range. Example: `80` or `8000-9000`
* `port_type` - (Optional) Port type, Available: local, remote
* `protocol` - (Required) Layer 4 protocol.  Available: tcp, udp, icmp, ip
## Attributes Reference

The following attributes are exported:

* `id` - The ID of firewall
* `server_id` - The ID of server
* `ip_address` -  The ip address of nic that will be used for apply firewall
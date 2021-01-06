---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_firewall_vpc"
sidebar_current: "docs-cmccloud-resource-firewall-vpc"
description: |-
  Provides a CMC Cloud Firewall VPC resource. This can be used to create, modify, and delete Firewall VPC.
---

# cmccloud\_firewall_vpc

Provides a CMC Cloud Firewall VPC resource. This can be used to create,
modify, and delete Firewall VPC. This firewall will be used for networks in VPC

## Example Create a Firewall VPC

```hcl
# Create a new Firewall for VPC
resource "cmccloud_firewall_vpc" "terra_firewall_1" {
    name = "terraform.firewall1"
    description = "vpc firewall create from terraform"
    vpc_id = cmccloud_vpc.terra_vpc.id
    inbound_rule {
        protocol         = "tcp"
        port_range       = "22"
        cidrs = ["192.168.1.0/24", "192.168.2.0/24"]
    } 
    outbound_rule {
        protocol         = "tcp"
        port_range       = "22"
        cidrs  = ["192.168.1.0/24", "192.168.2.0/24"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of firewall
* `description` - (Optional) The description for firewall
* `vpc_id` - (Required) The ID of VPC
* `inbound_rule` - (Optional) Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below
* `outbound_rule` - (Optional) Can be specified multiple times for each egress rule. Each egress block supports fields documented below.

The `inbound_rule` and `outbound_rule` block support the following: 

* `cidr` - (Required) CIDR Block: IPv4 or IPv6 CIDR
* `action` - (Optional) Action for this rule Available: allow, deny, default = allow
* `protocol` - (Required) Layer 4 protocol.  Available: tcp, udp, icmp
* `port_range` - (Optional) Port range. Example: `80` or `8000-9000`
* `cidrs` - The CIDR list to allow traffic from/to

## Attributes Reference

The following attributes are exported:

* `id` - The id of firwall
* `name` - The name of firewall
* `description` - The description for firewall
* `vpc_id` - The ID of VPC

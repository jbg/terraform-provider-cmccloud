---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_network"
sidebar_current: "docs-cmccloud-resource-network"
description: |-
  Provides a CMC Cloud Network resource. This can be used to create, modify, and delete network.
---

# cmccloud\_network

Provides a CMC Cloud network resource. This can be used to create,
modify, and delete network.

## Example Create Network

```hcl
# Create a new network
resource "cmccloud_network" "private_network1" {
    name = "terraform.private.network1"
    description = "network create from terraform"
    gateway = "10.10.0.1"
    netmask = "255.255.255.0"
    firewall_id = "allow"
    vpc_id = cmccloud_vpc.terra_vpc.id
    server_ids = [ cmccloud_server.terraform_server1.id, "cf967e0a-2784-4245-8b7e-583b89d7fb66" ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of network 
* `description` - (Optional) The description for network
* `gateway` - (Required) The gateway of the network 
* `netmask` - (Required) The netmask of the network.
* `vpc_id` - (Required) The VPC network belongs to
* `firewall_id` - (Required) Network Firewall ID associated for the network, avaiable value: allow,deny or Firewall ID
* `server_ids` - (Required) List of server ids that attach this network

## Attributes Reference

The following attributes are exported:

* `name` - The name of network 
* `description` - The description for network
* `gateway` - The gateway of the network 
* `netmask` - The netmask of the network.
* `vpc_id` - The VPC network belongs to
* `firewall_id` - Network Firewall ID associated for the network, avaiable value: allow,deny or Firewall ID
* `server_ids` - List of server ids that attach this network
* `cidr` - The network cidr
* `state` - The network state
* `type` - The network type

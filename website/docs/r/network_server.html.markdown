---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_network_server"
sidebar_current: "docs-cmccloud-resource-network-server"
description: |-
  Provides a CMC Cloud Server Network resource. This can be used to attach, detach network from server.
---

# cmccloud\_network\_server

Provides a CMC Cloud network resource. This can be used to attach, detach network from server.

## Example Create Server Network

```hcl
# Create a new server nic attach to a network
resource "cmccloud_network_server" "server_nic" {
    network_id = "e436e7fe-5506-4d6b-bbd3-1a0570a9c83b"
    server_id = "cf967e0a-2784-4245-8b7e-583b89d7fb16"
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required) The ID of network 
* `server_id` - (Required) The ID of server

## Attributes Reference

The following attributes are exported:

* `id` - The ID of created nic
* `network_id` - The ID of network 
* `server_id` - The ID of server

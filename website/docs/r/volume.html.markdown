---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_volume"
sidebar_current: "docs-cmccloud-resource-volume"
description: |-
  Provides a CMC Cloud Volume resource. This can be used to create, modify, and delete volumes.
---

# cmccloud\_volume

Provides a CMC Cloud Volume resource. This can be used to create,
modify, and delete volume.
## Example Usage

```hcl
# Create a new volume
resource "cmccloud_volume" "volume1" {
    name    = "vol1.ssd"
    region  = "HN"
    size    = 20
	type    = "ssd"
	server_id = "040e6918-409e-4346-8827-cd328e113f91"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `region` - (Required) The availability zone of the volume. Example: hn, hcm
* `size` - (Required) The size of the volume in GB
* `type` - (Required) The type of the volume: ssd or hdd.
* `server_id` - (Optional) - The id of server to attach


## Attributes Reference

The following attributes are exported:

* `name` - The name of the volume.
* `region` - The availability zone of the server. Example: hn, hcm
* `size` - The size of the volume in GB
* `type` - The type of the volume: ssd or hdd.
* `server_id` - The id of server to attach
* `state` - The state of volume
* `created` - Created time

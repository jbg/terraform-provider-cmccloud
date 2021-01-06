---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_server"
sidebar_current: "docs-cmccloud-resource-server"
description: |-
  Provides a CMC Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# cmccloud\_server

Provides a CMC Cloud Server resource. This can be used to create,
modify, and delete Server. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new Web Server
resource "cmccloud_server" "web" {
	name = "terraform.com"
    cpu = 1
    ram = 1
    root = 20
	gpu = 0
    image_type = "image"
    image_id = "472c8730-bc28-42d6-8680-151402be1169"
    region = "hn"
	auto_backup = true
	enable_private_network = false
	backup_schedule = "DAILY-03:00"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Server name.
* `cpu` - (Required) Number of CPU cores
* `ram` - (Required) Number of Memory in GB format
* `root` - (Required) Size of the root volume in GB format
* `gpu` - (Optional) Number of GPU card
* `image_type` - (Required) The type for create server root disk: image, snapshot, backup
* `image_id` - (Required) The ID of image - image ID, snapshot ID or backup ID 
* `region` - (Required) The availability zone of the server. Example: hn, hcm
* `auto_backup` - (Optional) Set to true to enable auto backup for root volume
* `enable_private_network` - (Optional) Set to true to enable private network
* `backup_schedule` - (Optional) When auto_backup is set to true, you can define the backup time in format: DAILY-HH:mm, WEEKLY-HH:mm:w, MONTHLY-HH:mm:dd (HH: hours, mm: minute, w: day of week - 0 (for Sunday) through 6 (for Saturday) and dd: day of month)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Server
* `name` - The Server name.
* `cpu` - Number of CPU cores
* `ram` - Number of Memory in GB format
* `root` - Size of the root volume in GB format
* `gpu` - Number of GPU card
* `image_type` - The type for create server root disk: image, snapshot, backup
* `image_id` - The ID of OS - image ID, snapshot ID or backup ID 
* `region` - The availability zone of the server. Example: hn, hcm
* `auto_backup` - state of auto backup
* `enable_private_network` - Enable or disable private network
* `backup_schedule` Backup time
* `created` Created time
* `image_name` Name of the image
* `state` State of server
* `main_ip_address` Main ip address of server
---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_snapshot"
sidebar_current: "docs-cmccloud-resource-snapshot"
description: |-
  Provides a CMC Cloud Volume Snapshot resource. This can be used to create and delete volume snapshot.
---

# cmccloud\_snapshot

Provides a CMC Cloud Volume Snapshot resource. This can be used to create,
and delete volume snapshot.
## Example Usage

```hcl
# Create a new snapshot
resource "cmccloud_snapshot" "snapshot1" {
    name = "snapshot1"
    volume_id = "cc3cce08-f514-4186-8e19-dbc38fb8f6ef"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the snapshot.
* `volume_id` - (Optional) The ID of volume will be take snapshot. Required if server_id is missing
* `server_id` - (Optional) The ID of server will be take snapshot. Required if volume_id is missing


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the snapshot
* `name`- The name of the snapshot
* `volume_id` - The ID of volume
* `server_id` - The ID of server
* `size` - The size of snapshot
* `state`- State of the snapshot
* `created` - Created time
* `volume_name` - The name of volume 
* `server_name` - The name of server

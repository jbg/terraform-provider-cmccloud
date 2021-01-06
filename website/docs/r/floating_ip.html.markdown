---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_floating_ip"
sidebar_current: "docs-cmccloud-resource-floating-ip"
description: |-
  Provides a CMC Cloud Floating IP resource. This can be used to create and delete floating ip.
---

# cmccloud\_floating\_ip

Provides a CMC Cloud Floating IP resource. This can be used to create,
and delete floating ip
## Example Usage

```hcl
# Create a new Floating IP
resource "cmccloud_floating_ip" "terra_floatip_1"{
    vpc_id = "040e6918-409e-4346-8827-cd322e113f23"
}
```

## Argument Reference
 
* `vpc_id` - (Required) The ID of the VPC the ip belongs to

## Attributes Reference

The following attributes are exported:

* `vpc_id` - The ID of the VPC the ip belongs to
* `ip_address` - Public IP address
* `is_source_nat` - true if the IP address is a source nat address, false otherwise
* `is_static_nat` - true if this ip is for static nat, false otherwise
* `state` - State of the ip address. Can be: Allocating, Allocated and Releasing

---
layout: "cmccloud"
page_title: "CMC Cloud: cmccloud_vpc"
sidebar_current: "docs-cmccloud-resource-vpc"
description: |-
  Provides a CMC Cloud VPC resource. This can be used to create and delete vpc.
---

# cmccloud\_vpc

Provides a CMC Cloud VPC resource. This can be used to create,
and delete vpc
## Example Usage

```hcl
# Create a new VPC
resource "cmccloud_vpc" "vpc1" {
    name = "terraform.vpc"
    description = "VPC for hn"
    region = "hn"
    cidr = "10.10.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPC.
* `description` - (Required) The description text of the VPC
* `region` - (Required) The availability zone of the vpc. Example: hn, hcm
* `cidr` - (Required) The cidr of the VPC. All VPC guest networks' cidrs should be within this CIDR

## Attributes Reference

The following attributes are exported:

* `name` - The name of the VPC.
* `description` - The description text of the VPC
* `region` - The availability zone of the vpc. Example: hn, hcm
* `cidr` - The cidr of the VPC. All VPC guest networks' cidrs should be within this CIDR
* `state` - The state of VPC

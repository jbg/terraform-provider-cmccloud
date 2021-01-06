---
layout: "cmccloud"
page_title: "Provider: CMC Cloud"
sidebar_current: "docs-cmccloud-index"
description: |-
  The CMC Cloud provider is used to interact with the resources supported by CMC Cloud. The provider needs to be configured with the proper credentials before it can be used.
---


 CMC Cloud Provider

The CMC Cloud provider is used to interact with the
resources supported by CMC Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the CMC Cloud Provider
provider "cmccloud" {
    api_endpoint = "https://api.cloud.cmctelecom.vn/ver2"
    api_key = "password"
}

# Create a server
resource "cmccloud_server" "test.server" {
  # ...
}
```

## Argument Reference

The following arguments are supported:
 
* `api_endpoint` - (Optional) This can be used to override the base URL for
  CMC Cloud API requests (Defaults to the value of the `CMC_CLOUD_API_ENDPOINT`
  environment variable or `https://api.cloud.cmctelecom.vn/ver2` if unset).
  
* `api_key` - (Required) This is your api key get from https://portal.cloud.cmctelecom.vn to authenticate with CMC Cloud.  Alternatively, this can also be specified using environment 
  variables ordered by precedence: `CMC_CLOUD_API_KEY`
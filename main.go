package main

import (
	"terraform-provider-cmccloud/cmccloud"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cmccloud.Provider})
}

package cmccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for CMC Cloud.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL use for the CMC Cloud API",
				DefaultFunc: schema.EnvDefaultFunc("CMC_CLOUD_API_ENDPOINT", "https://api.cloud.cmctelecom.vn"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API key get from account settings in https://portal.cloud.cmctelecom.vn",
				DefaultFunc: schema.EnvDefaultFunc("CMC_CLOUD_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cmccloud_server":          resourceCMCCloudServer(),
			"cmccloud_volume":          resourceCMCCloudVolume(),
			"cmccloud_snapshot":        resourceCMCCloudSnapshot(),
			"cmccloud_vpc":             resourceCMCCloudVPC(),
			"cmccloud_network":         resourceCMCCloudNetwork(),
			"cmccloud_floating_ip":     resourceCMCCloudFloatingIP(),
			"cmccloud_firewall_vpc":    resourceCMCCloudFirewallVPC(),
			"cmccloud_firewall_direct": resourceCMCCloudFirewallDirect(),
			"cmccloud_network_server":  resourceCMCCloudNetworkServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"cmccloud_images": datasourceCMCCloudImages(),
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	} /**/
	return p
}

/**/
func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		APIEndpoint:      d.Get("api_endpoint").(string),
		APIKey:           d.Get("api_key").(string),
		TerraformVersion: terraformVersion,
	}
	return config.Client()
}

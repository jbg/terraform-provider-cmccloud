package cmccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceCMCCloudImages() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceCMCCloudImageRead,
		Schema: imagesSchema(),
	}
}

func dataSourceCMCCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	osImages, err := client.Image.List(context.Background())
	if err != nil {
		return err
	}
	d.Set("images", osImages)
	return nil
}

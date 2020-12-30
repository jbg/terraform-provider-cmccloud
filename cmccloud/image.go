package cmccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Id of the image",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "OS Distribution",
			Computed:    true,
		},
		"bits": {
			Type:        schema.TypeString,
			Description: "OS Bit",
			Computed:    true,
		},
	}
}

func imagesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"images": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: imageSchema(),
			},
		},
	}
}

package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func serverSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		}, "label": {
			Type:     schema.TypeString,
			Required: true,
		}, "location": {
			Type:     schema.TypeString,
			Required: true,
		}, "product_code": {
			Type:     schema.TypeString,
			Required: true,
		}, "image": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

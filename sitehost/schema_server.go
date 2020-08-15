package sitehost

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/internal/api"
)

const (
	nameKey        = "name"
	labelKey       = "label"
	locationKey    = "location"
	productCodeKey = "product_code"
	imageKey       = "image"
	passwordKey    = "password"
	ipsKey         = "ips"
)

var schemaServer = map[string]*schema.Schema{
	nameKey: {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}, labelKey: {
		Type:     schema.TypeString,
		Required: true,
	}, locationKey: {
		Type:     schema.TypeString,
		Required: true,
	}, productCodeKey: {
		Type:     schema.TypeString,
		Required: true,
	}, imageKey: {
		Type:     schema.TypeString,
		Required: true,
	},
}

func unmarshalServer(s api.Server) map[string]interface{} {
	return map[string]interface{}{
		nameKey:        s.Name,
		labelKey:       s.Label,
		locationKey:    s.Location,
		productCodeKey: s.ProductCode,
	}
}

func unmarshalServers(s []api.Server) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(s))
	for _, s := range s {
		m = append(m, unmarshalServer(s))
	}
	return m
}

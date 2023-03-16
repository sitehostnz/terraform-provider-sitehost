package dns

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sitehostnz/gosh/pkg/utils"
)

// resourceZoneSchema is the schema with values for a DNS zone resource.
var resourceZoneSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The domain name",
	},
}

// resourceRecordSchema is the schema with values for a DNS record resource.
var resourceRecordSchema = map[string]*schema.Schema{
	"domain": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The base domain",
	},

	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     false,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The subdomain",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			domain := d.Get("domain").(string)

			oldValue = utils.ConstructFqdn(oldValue, domain)
			newValue = utils.ConstructFqdn(newValue, domain)

			return newValue == oldValue
		},
	},

	"type": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"A",
			"AAAA",
			"CAA",
			"CNAME",
			"MX",
			"NS",
			"TXT",
			"SRV",
		}, false),
		Description: "The record type",
	},

	"priority": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description:  "The priority type",
		Default:      0,
	},

	"ttl": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},

	"content": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, oldRecord, newRecord string, d *schema.ResourceData) bool {
			return strings.TrimSuffix(oldRecord, ".") == strings.TrimSuffix(newRecord, ".")
		},
	},

	"fqdn": {
		Type:     schema.TypeString,
		Computed: true,
	},

	"change_date": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

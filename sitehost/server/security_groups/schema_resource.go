package security_groups

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceSchema is the schema with values for a Server resource.
var resourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The name of the Security Group.",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label for the Security Group.",
	},
	"rules": {
		Type: schema.TypeMap,
		Elem: map[string]*schema.Schema{
			"in": {
				Type: schema.TypeList,
				Elem: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"action": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ACCEPT",
							"DROP",
							"REJECT",
						}, false),
						Description: "Whether this inbound rule is enabled or not.",
					},
					"protocol": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"tcp",
							"udp",
						}, false),
						Description: "Whether this inbound rule is enabled or not.",
					},
					"src_ip": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"dest_port": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
				},
			},
			"out": {
				Type: schema.TypeList,
				Elem: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"action": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ACCEPT",
							"DROP",
							"REJECT",
						}, false),
						Description: "Whether this inbound rule is enabled or not.",
					},
					"protocol": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"tcp",
							"udp",
						}, false),
						Description: "Whether this inbound rule is enabled or not.",
					},
					"dest_ip": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"dest_port": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
				},
			},
		},
		Optional:    true,
		Description: "The security groups which this server uses.",
	},
}

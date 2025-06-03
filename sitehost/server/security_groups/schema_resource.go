package security_groups

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceSchema is the schema with values for a Security Group resource.
var resourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The name of the Security Group.",
	},
	"label": {
		Type:        schema.TypeString,
		Optional:    true,
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
						Optional:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"action": {
						Type:     schema.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ACCEPT",
							"DROP",
							"REJECT",
						}, false),
						Description: "The action for this inbound rule. The following values are accepted: ACCEPT, DROP, REJECT.",
					},
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"tcp",
							"udp",
						}, false),
						Description: "The protocol for this inbound rule. The following values are accepted: tcp, udp.",
					},
					"src_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The source IP address for this inbound rule. This can either be a standalone IP or CIDR range.",
					},
					"dest_port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The destination port for this inbound rule.",
					},
				},
			},
			"out": {
				Type: schema.TypeList,
				Elem: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether this inbound rule is enabled or not.",
					},
					"action": {
						Type:     schema.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ACCEPT",
							"DROP",
							"REJECT",
						}, false),
						Description: "The action for this inbound rule. The following values are accepted: ACCEPT, DROP, REJECT.",
					},
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"tcp",
							"udp",
						}, false),
						Description: "The protocol for this inbound rule. The following values are accepted: tcp, udp.",
					},
					"dest_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The destination IP address for this outbound rule. This can either be a standalone IP or CIDR range.",
					},
					"dest_port": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The destination port for this inbound rule.",
					},
				},
			},
		},
		Required:    true,
		Description: "The rules which the security group follows.",
	},
}

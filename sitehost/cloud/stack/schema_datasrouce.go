package stack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// // resourceSchema is the schema with values for a Server resource.
// DataSource returns a schema with the function to read Server resource.
func stackDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The client id",
		},
		"server_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The server id",
		},
		"server_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The server name",
		},
		"server_label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The server name",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The stack name",
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The stack name",
		},
	}
}

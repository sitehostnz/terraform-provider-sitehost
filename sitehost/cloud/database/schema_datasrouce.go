package database

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// // resourceSchema is the schema with values for a Server resource.
// DataSource returns a schema with the function to read Server resource.
func stackDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"databases": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "A list of databases",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The client id",
					},
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
						Computed:    true,
						Description: "The server name",
					},
					"server_label": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The server label",
					},
					"server_ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The server ip address",
					},
					"mysql_host": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The the mysql host",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The db name",
					},
					"container": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The container name",
					},
				},
			},
		},
	}
}

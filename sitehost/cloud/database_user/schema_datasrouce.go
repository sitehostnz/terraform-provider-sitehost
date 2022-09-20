package database_user

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
			Description: "A list of database users",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{},
			},
		},
	}
}

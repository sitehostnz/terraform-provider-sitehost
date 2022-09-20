package database_user

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/hashcode"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"strconv"
)

func DataSource() *schema.Resource {
	recordSchema := stackDataSourceSchema()

	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client
	fetchedDatabaseUsers, err := client.DatabaseUsers.List(ctx)

	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	hashBase, _ := json.Marshal(fetchedDatabaseUsers)
	hashString := strconv.Itoa(hashcode.String(string(hashBase)))
	d.SetId(hashString)

	databases := make([]map[string]interface{}, len(*fetchedDatabaseUsers))
	for i, _ := range *fetchedDatabaseUsers {
		databases[i] = map[string]interface{}{
			//"id":         db.Id,
			//"client_id":  db.ClientId,
			//"mysql_host": db.MysqlHost,
			//"container":  db.Container,
			//"name":       db.DbName,
			//// server is nested hell times...
			//"server_id":    db.ServerId,
			//"server_ip":    db.ServerIp,
			//"server_name":  db.ServerName,
			//"server_label": db.ServerLabel,
		}
	}
	d.Set("users", databases)

	return nil
}

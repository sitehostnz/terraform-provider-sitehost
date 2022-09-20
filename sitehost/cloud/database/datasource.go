package database

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
	fetchedDatabases, err := client.Database.List(ctx)

	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	hashBase, _ := json.Marshal(fetchedDatabases)
	hashString := strconv.Itoa(hashcode.String(string(hashBase)))
	d.SetId(hashString)

	databases := make([]map[string]interface{}, len(*fetchedDatabases))
	for i, db := range *fetchedDatabases {
		databases[i] = map[string]interface{}{
			"id":         db.Id,
			"client_id":  db.ClientId,
			"mysql_host": db.MysqlHost,
			"container":  db.Container,
			"name":       db.DbName,
			// server is nested hell times...
			"server_id":    db.ServerId,
			"server_ip":    db.ServerIp,
			"server_name":  db.ServerName,
			"server_label": db.ServerLabel,
		}
	}
	d.Set("databases", databases)

	return nil
}

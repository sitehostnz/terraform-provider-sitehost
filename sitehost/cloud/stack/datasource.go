package stack

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
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
	stack, err := client.Stack.Get(ctx, &gosh.Stack{
		ClientID: client.ClientID,
		Server:   d.Get("server_name").(string),
		Name:     d.Get("name").(string),
	})
	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	d.SetId(stack.Name)
	d.Set("server_id", stack.ServerID)
	d.Set("label", stack.Label)
	d.Set("server_name", stack.Server)
	d.Set("server_label", stack.ServerLabel)

	return nil
}

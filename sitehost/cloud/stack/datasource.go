package stack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
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
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := stack.New(conf.Client)

	stack, err := client.Get(ctx, stack.GetRequest{
		ServerName: d.Get("server_name").(string),
		Name:       d.Get("name").(string),
	})
	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	d.SetId(stack.Stack.Name)
	d.Set("server_id", stack.Stack.ServerID)
	d.Set("label", stack.Stack.Label)
	d.Set("server_name", stack.Stack.Server)
	d.Set("server_label", stack.Stack.ServerLabel)

	return nil
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServer() *schema.Resource {
	recordSchema := serverSchema()

	return &schema.Resource{
		ReadContext: dataSourceServerRead,
		Schema:      recordSchema,
	}
}

func dataSourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

package sitehost

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/internal/api"
)

var dataSourceServers = schema.Resource{
	ReadContext: crud(dataSourceServersRead),
	Schema: map[string]*schema.Schema{
		"servers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemaServer,
			},
		},
	},
}

func dataSourceServersRead(ctx context.Context, c api.Client, d *schema.ResourceData) error {
	ss, err := c.Server().List(ctx)
	if err != nil {
		return err
	}
	if err = d.Set("servers", unmarshalServers(ss)); err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

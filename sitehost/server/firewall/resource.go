// Package firewall provides the functions to create a Firewall resource via SiteHost API.
package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server/firewall"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for server.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: updateResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// readResource is a function to read the firewall of a server.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	resp, err := client.Get(ctx, firewall.GetRequest{
		ServerName: d.Get("server").(string),
	})
	if err != nil {
		return diag.Errorf("Error reading server: %s", err)
	}

	groups := make([]string, len(resp.Return))
	for i, group := range resp.Return {
		groups[i] = group.Group
	}
	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update the firewall of a server.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	serverName, ok := d.Get("server").(string)
	if !ok {
		return diag.Errorf("failed to convert server name to string")
	}
	list, ok := d.Get("groups").([]interface{})
	if !ok {
		return diag.Errorf("failed to convert groups to []interface{}")
	}
	groups := make([]string, 0)
	for _, v := range list {
		groups = append(groups, v.(string))
	}

	diags := updateFirewallGroups(client, serverName, groups)
	if diags != nil {
		return diags
	}

	d.SetId(serverName)
	return nil
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	serverName, ok := d.Get("server").(string)
	if !ok {
		return diag.Errorf("failed to convert server name to string")
	}

	diags := updateFirewallGroups(client, serverName, []string{})
	if diags != nil {
		return diags
	}

	// Clear the ID when the resource is deleted
	d.SetId("")
	return nil
}

func updateFirewallGroups(client *firewall.Client, serverName string, groups []string) diag.Diagnostics {
	res, err := client.Update(context.Background(), firewall.UpdateRequest{
		ServerName:     serverName,
		SecurityGroups: groups,
	})
	if err != nil {
		return diag.Errorf("Error updating server: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error updating server: %s", res.Msg)
	}

	return nil
}

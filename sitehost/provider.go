// Package sitehost provides the functions to create a Terraform Provider to create resources via SiteHost API.
package sitehost

import (
	"context"
	"fmt"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/api_info"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/domain"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/domain_record"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/cloud/stack"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/dns"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/info"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/server"
)

// New returns a schema.Provider for SiteHost.
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SH_CLIENT_ID", nil),
					Description: "The client identifier that allows you access to your SiteHost account.",
				}, "api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SH_APIKEY", nil),
					Description: "The API Key that allows you access to your SiteHost account.",
					Sensitive:   true,
				}, "api_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The HTTP(S) API address of the SiteHost API to use.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"sitehost_server": server.DataSource(),
				"sitehost_api":    info.DataSource(),
				"sitehost_stack":  stack.DataSource(),
				// "sitehost_stack_database": database.DataSource(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"sitehost_server":     server.Resource(),
				"sitehost_dns_zone":   dns.ZoneResource(),
				"sitehost_dns_record": dns.RecordResource(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
			return configure(ctx, version, d)
		}

		return p
	}
}

// configure returns the Config with connection data.
func configure(_ context.Context, version string, d *schema.ResourceData) (any, diag.Diagnostics) {
	config := &helper.Config{
		APIKey:           fmt.Sprint(d.Get("api_key")),
		ClientID:         fmt.Sprint(d.Get("client_id")),
		APIEndpoint:      fmt.Sprint(d.Get("api_endpoint")),
		TerraformVersion: version,
	}

	return config.Client()
}

// Package sitehost provides the functions to create a Terraform Provider to create resources via SiteHost API.
package sitehost

import (
	"context"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/api"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/domain"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/domain_record"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
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
					Description: "The client identifier that allows you access to your SiteHost api.",
				}, "api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SH_APIKEY", nil),
					Description: "The API Key that allows you access to your SiteHost api.",
					Sensitive:   true,
				}, "api_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The HTTP(S) API address of the SiteHost API to use.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"sitehost_api":    api.DataSource(),
				"sitehost_server": server.DataSource(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"sitehost_server":        server.Resource(),
				"sitehost_domain":        domain.Resource(),
				"sitehost_domain_record": domain_record.Resource(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return configure(ctx, version, d)
		}

		return p
	}
}

// configure returns the Config with connection data.
func configure(_ context.Context, version string, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := &helper.Config{
		APIKey:           d.Get("api_key").(string),
		ClientID:         d.Get("client_id").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		TerraformVersion: version,
	}

	return config.Client()
}

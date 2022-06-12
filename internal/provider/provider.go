package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// New return a schema.Provider for SiteHost
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SH_CLIENT_ID", nil),
					Description: "client identifier",
				}, "api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SH_APIKEY", nil),
					Description: "api authentication key",
					Sensitive:   true,
				}, "api_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "url prefix of the api server",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"sitehost_server": dataSourceServer(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"sitehost_server": resourceServer(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return configure(ctx, version, d)
		}

		return p
	}
}

// configure return the Config with connection data
func configure(_ context.Context, version string, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		APIKey:           d.Get("api_key").(string),
		ClientID:         d.Get("client_id").(string),
		APIEnpoint:       d.Get("api_endpoint").(string),
		TerraformVersion: version,
	}

	return config.Client()
}

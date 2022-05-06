package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sitehostnz/gosh"
	"log"
	"net/url"
	"time"
)

type Config struct {
	APIKey           string
	ClientID         string
	APIEnpoint       string
	TerraformVersion string
}

type CombinedConfig struct {
	client   *gosh.Client
	apiKey   string
	clientID string
}

func (c *CombinedConfig) goshClient() *gosh.Client { return c.client }

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	goshClient := gosh.NewClient(c.APIKey, c.ClientID)

	goshClient.UserAgent = fmt.Sprintf("Terraform/%s", c.TerraformVersion)

	if c.APIEnpoint != "" {
		apiURL, err := url.Parse(c.APIEnpoint)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		goshClient.BaseURL = apiURL
		log.Printf("[INFO] Sitehost Client configured for URL: %s", goshClient.BaseURL.String())
	}

	return &CombinedConfig{
		client:   goshClient,
		apiKey:   c.APIKey,
		clientID: c.ClientID,
	}, nil
}

func waitForAction(client *gosh.Client, jobID string) error {
	var (
		pending   = "Pending"
		target    = "Completed"
		ctx       = context.Background()
		refreshFn = func() (result interface{}, state string, err error) {
			j, err := client.Jobs.Get(ctx, jobID)
			if err != nil {
				return nil, "", err
			}
			if j.Return.State == "Failed" {
				return j, "Failed", nil
			}
			if j.Return.State == target {
				return j, target, nil
			}

			return j, pending, nil
		}
	)

	_, err := (&resource.StateChangeConf{
		Pending: []string{pending},
		Refresh: refreshFn,
		Target:  []string{target},

		Delay:          10 * time.Second,
		Timeout:        60 * time.Minute,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}).WaitForStateContext(ctx)

	return err
}

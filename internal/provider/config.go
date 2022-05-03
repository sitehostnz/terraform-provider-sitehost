package provider

import (
	"fmt"
	"github.com/sitehostnz/gosh"
)

type Config struct {
	APIKey           string
	ClientID         string
	TerraformVersion string
}

type CombinedConfig struct {
	client   *gosh.Client
	apiKey   string
	clientID string
}

func (c *CombinedConfig) goshClient() *gosh.Client { return c.client }

func (c *Config) Client() *CombinedConfig {
	goshClient := gosh.NewClient(c.APIKey, c.ClientID)

	goshClient.UserAgent = fmt.Sprintf("Terraform/%s", c.TerraformVersion)

	return &CombinedConfig{
		client:   goshClient,
		apiKey:   c.APIKey,
		clientID: c.ClientID,
	}
}

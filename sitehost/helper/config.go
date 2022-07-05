package helper

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sitehostnz/gosh"
)

const (
	JobStatusPending         = "Pending"
	JobStatusCompleted       = "Completed"
	JobStatusFailed          = "Failed"
	JobRequestDelay          = 10 * time.Second
	JobRequestTimeout        = 60 * time.Minute
	JobRequestMinTimeout     = 3 * time.Second
	JobRequestNotFoundChecks = 60
)

type Config struct {
	APIKey           string
	ClientID         string
	APIEndpoint      string
	TerraformVersion string
}

type CombinedConfig struct {
	Client *gosh.Client
	Config *Config
}

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	client := gosh.NewClient(c.APIKey, c.ClientID)

	client.UserAgent = "Terraform/" + c.TerraformVersion

	if c.APIEndpoint != "" {
		apiURL, err := url.Parse(c.APIEndpoint)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		client.BaseURL = apiURL
		log.Printf("[INFO] SiteHost Client configured for URL: %s", client.BaseURL.String())
	}

	return &CombinedConfig{
		Client: client,
		Config: c,
	}, nil
}

func WaitForAction(client *gosh.Client, jobID string) error {
	var (
		pending   = JobStatusPending
		target    = JobStatusCompleted
		ctx       = context.Background()
		refreshFn = func() (result interface{}, state string, err error) {
			j, err := client.Jobs.Get(ctx, jobID)
			if err != nil {
				return nil, "", err
			}
			if j.Return.State == JobStatusFailed {
				return j, JobStatusFailed, nil
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

		Delay:          JobRequestDelay,
		Timeout:        JobRequestTimeout,
		MinTimeout:     JobRequestMinTimeout,
		NotFoundChecks: JobRequestNotFoundChecks,
	}).WaitForStateContext(ctx)

	return err
}

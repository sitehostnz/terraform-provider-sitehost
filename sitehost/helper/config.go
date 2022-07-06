// Package helper provides the functions to work with SiteHost API.
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
	// JobStatusPending the name status for a pending job.
	JobStatusPending = "Pending"
	// JobStatusCompleted the name status for a completed job.
	JobStatusCompleted = "Completed"
	// JobStatusFailed  the name status for a failed job.
	JobStatusFailed = "Failed"
	// JobRequestDelay the time wait to send a new request to check the job status.
	JobRequestDelay = 10 * time.Second
	// JobRequestTimeout the time to wait before timeout.
	JobRequestTimeout = 60 * time.Minute
	// JobRequestMinTimeout smallest time to wait before refreshes.
	JobRequestMinTimeout = 3 * time.Second
	// JobRequestNotFoundChecks number of times to allow not found.
	JobRequestNotFoundChecks = 60
)

// Config is a wrapper to save the configuration connection from terraform.
type Config struct {
	APIKey           string
	ClientID         string
	APIEndpoint      string
	TerraformVersion string
}

// CombinedConfig is s struct with API wrapper and the Config.
type CombinedConfig struct {
	Client *gosh.Client
	Config *Config
}

// Client returns a new CombinedConfig instance.
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

// WaitForAction function to check the Job status in a refresh function.
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

// Package server provides the functions to create a Server resource via SiteHost API.
package domain

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"log"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new Server resource.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	domain, err := client.Domain.Create(
		ctx, &gosh.Domain{
			ClientID: client.ClientID,
			Name:     d.Get("name").(string),
		})

	if err != nil {
		return diag.Errorf("Error creating domain: %s", err)
	}

	if domain == nil {
		return diag.Errorf("Failed retrieving domain: %s", err)
	}

	d.SetId(domain.Name)

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	domain, err := client.Domain.Get(
		context.Background(),
		&gosh.Domain{
			ClientID: client.ClientID,
			Name:     d.Id(),
		},
	)

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if domain == nil {
		return diag.Errorf("Failed retrieving domain: %s", err)
	}

	d.SetId(domain.Name)
	d.Set("name", domain.Name)

	log.Printf("[INFO] Domain Name: %s", domain.Name)

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client
	//
	_, err := client.Domain.Delete(context.Background(), &gosh.Domain{
		Name:     d.Id(),
		ClientID: client.ClientID,
	})

	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	return nil
}

// Package server provides the functions to create a Server resource via SiteHost API.
package server

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,

		Schema: resourceSchema,
	}
}

// createResource is a function to create a new Server resource.
func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := conf.Client
	keys, ok := d.Get("ssh_keys").([]any)
	if !ok {
		return diag.Errorf("failed to convert ssh keys object")
	}

	sshKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		sshKeys = append(sshKeys, fmt.Sprint(key))
	}

	opts := &gosh.ServerCreateRequest{
		Label:       fmt.Sprint(d.Get("label")),
		Location:    fmt.Sprint(d.Get("location")),
		ProductCode: fmt.Sprint(d.Get("product_code")),
		Image:       fmt.Sprint(d.Get("image")),
		Params: gosh.ParamsOptions{
			SSHKeys: sshKeys,
		},
	}

	res, err := client.Servers.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating server: %s", err)
	}

	// Set data
	d.SetId(res.Return.Name)
	if err = d.Set("name", res.Return.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("password", res.Return.Password); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("ips", res.Return.Ips); err != nil {
		return diag.FromErr(err)
	}

	// wait for "Completed" status
	if err = helper.WaitForAction(client, res.Return.JobID); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Server Name: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := conf.Client

	server, err := client.Servers.Get(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving server: %s", err)
	}

	if err = setServerAttributes(d, server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update a new Server resource.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := conf.Client

	if d.HasChange("product_code") {
		if err := client.Servers.Upgrade(context.Background(), &gosh.ServerUpgradeRequest{Name: d.Id(), Plan: fmt.Sprint(d.Get("product_code"))}); err != nil {
			return diag.FromErr(err)
		}

		resp, err := client.Servers.CommitChanges(context.Background(), d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if err = helper.WaitForAction(client, resp.Return.JobID); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	if d.HasChange("label") {
		if err := client.Servers.Update(context.Background(), &gosh.ServerUpdateRequest{Name: d.Id(), Label: fmt.Sprint(d.Get("label"))}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	return readResource(ctx, d, meta)
}

// deleteResource is a function to delete a new Server resource.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := conf.Client

	resp, err := client.Servers.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	if err = helper.WaitForAction(client, resp.Return.JobID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setServerAttributes is a function to set data to a Server resource.
func setServerAttributes(d *schema.ResourceData, server *gosh.Server) error {
	return d.Set("name", server.Name)
}

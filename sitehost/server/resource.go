// Package server provides the functions to create a Server resource via SiteHost API.
package server

import (
	"context"
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

// createResource function to create a new Server resource.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client
	keys := d.Get("ssh_keys").([]interface{})

	sshKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		sshKeys = append(sshKeys, key.(string))
	}

	opts := &gosh.ServerCreateRequest{
		Label:       d.Get("label").(string),
		Location:    d.Get("location").(string),
		ProductCode: d.Get("product_code").(string),
		Image:       d.Get("image").(string),
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
	err = d.Set("name", res.Return.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("password", res.Return.Password)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("ips", res.Return.Ips)
	if err != nil {
		return diag.FromErr(err)
	}

	// wait for "Completed" status
	err = helper.WaitForAction(client, res.Return.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Server Name: %s", d.Id())

	return nil
}

// readResource function to read a new Server resource.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	server, err := client.Servers.Get(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving server: %s", err)
	}

	err = setServerAttributes(d, server)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource function to update a new Server resource.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	if d.HasChange("product_code") {
		err := client.Servers.Upgrade(context.Background(), &gosh.ServerUpgradeRequest{Name: d.Id(), Plan: d.Get("product_code").(string)})
		if err != nil {
			return diag.FromErr(err)
		}

		resp, err := client.Servers.CommitChanges(context.Background(), d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		err = helper.WaitForAction(client, resp.Return.JobID)
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	if d.HasChange("label") {
		err := client.Servers.Update(context.Background(), &gosh.ServerUpdateRequest{Name: d.Id(), Label: d.Get("label").(string)})
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	return readResource(ctx, d, meta)
}

// deleteResource function to delete a new Server resource.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	resp, err := client.Servers.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	err = helper.WaitForAction(client, resp.Return.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setServerAttributes function to set data to a Server resource.
func setServerAttributes(d *schema.ResourceData, server *gosh.Server) error {
	err := d.Set("name", server.Name)
	if err != nil {
		return err
	}
	return nil
}

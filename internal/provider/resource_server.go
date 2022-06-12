package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			}, "password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			}, "ips": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			}, "label": {
				Type:     schema.TypeString,
				Required: true,
			}, "location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			}, "product_code": {
				Type:     schema.TypeString,
				Required: true,
			}, "image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			}, "ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).goshClient()

	keys := d.Get("ssh_keys").([]interface{})
	var sshKeys []string

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
	d.Set("name", res.Return.Name)
	d.Set("password", res.Return.Password)
	d.Set("ips", res.Return.Ips)

	// wait for "Completed" status
	err = waitForAction(client, res.Return.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Server Name: %s", d.Id())

	return nil
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).goshClient()

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

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).goshClient()

	if d.HasChange("product_code") {
		err := client.Servers.Upgrade(context.Background(), &gosh.ServerUpgradeRequest{Name: d.Id(), Plan: d.Get("product_code").(string)})
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Servers.CommitChanges(context.Background(), d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForAction(client, resp.Return.JobID)
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

	return resourceServerRead(ctx, d, meta)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).goshClient()

	resp, err := client.Servers.Delete(context.Background(), d.Id())

	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	err = waitForAction(client, resp.Return.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setServerAttributes(d *schema.ResourceData, server *gosh.Server) error {
	d.Set("name", server.Name)
	return nil
}

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
	"log"
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
				ForceNew: true,
			}, "image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).goshClient()

	opts := &gosh.ServerCreateRequest{
		Label:       d.Get("label").(string),
		Location:    d.Get("location").(string),
		ProductCode: d.Get("product_code").(string),
		Image:       d.Get("image").(string),
	}

	res, err := client.Servers.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating server: %s", err)
	}

	d.SetId(res.Return.Name)
	d.Set("password", res.Return.Password)

	err = waitForAction(client, res.Return.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Server Name: %s", d.Id())

	return nil
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	panic(nil)
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	panic(nil)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	panic(nil)
}

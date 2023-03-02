// Package dns provides the functions to create/get dns zones/records resource via SiteHost API.
package dns

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/dns"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// ZoneResource returns a schema with the operations for DNS resource.
func ZoneResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createZoneResource,
		ReadContext:   readZoneResource,
		DeleteContext: deleteZoneResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceZoneSchema,
	}
}

// createResource is a function to create a new DNS resource.
func createZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	domain := fmt.Sprintf("%v", d.Get("name"))

	// The response don't have the domain name, so we need to get it from the request.
	resp, err := client.CreateZone(ctx, dns.CreateZoneRequest{DomainName: domain})
	if err != nil {
		return diag.Errorf("Error creating domain: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error creating domain: %s", resp.Msg)
	}

	d.SetId(domain)

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// readResource is a function to read a new DNS resource.
func readZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	response, err := client.GetZone(ctx, dns.GetZoneRequest{DomainName: fmt.Sprint(d.Id())})
	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if !response.Status {
		return diag.Errorf("Error retrieving domain: %s", response.Msg)
	}

	// iterate over the domains to find the one we are looking for.
	for _, domain := range response.Return {
		if domain.Name == fmt.Sprint(d.Id()) {
			if err := d.Set("name", domain.Name); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	return diag.Errorf("Error finding the domain")
}

// deleteResource is a function to delete a new DNS resource.
func deleteZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.DeleteZone(ctx, dns.DeleteZoneRequest{DomainName: d.Id()})
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting server: %s", resp.Msg)
	}

	return nil
}

// RecordResource returns a schema with the operations for DNS resource.
func RecordResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createRecordResource,
		ReadContext:   readRecordResource,
		UpdateContext: updateRecordResource,
		DeleteContext: deleteRecordResource,
		Importer: &schema.ResourceImporter{
			StateContext: importRecordResource,
		},
		Schema: resourceRecordSchema,
	}
}

// createResource is a function to create a new Server resource.
func createRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	domain := fmt.Sprintf("%v", d.Get("domain"))
	name := fmt.Sprintf("%v", d.Get("name"))
	domainRecord := helper.ConstructFqdn(name, domain)

	client := dns.New(conf.Client)
	resp, err := client.AddRecord(ctx, dns.AddRecordRequest{
		Domain:   domain,
		Type:     fmt.Sprintf("%v", d.Get("type")),
		Name:     domainRecord,
		Content:  fmt.Sprintf("%v", d.Get("record")),
		Priority: fmt.Sprintf("%v", d.Get("priority")),
	})
	if err != nil {
		return diag.Errorf("Error creating DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error creating DNS record: %s", resp.Msg)
	}

	log.Printf("[INFO] Domain Record: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	domain := fmt.Sprintf("%v", d.Get("domain"))
	resp, err := client.ListRecords(ctx, dns.ListRecordsRequest{
		Domain: domain,
	})
	if err != nil {
		return diag.Errorf("Error retrieving DNS zone: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error creating DNS zone: %s", resp.Msg)
	}

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.DeleteRecord(ctx, dns.DeleteRecordRequest{
		Domain:   fmt.Sprintf("%v", d.Get("domain")),
		RecordID: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error deleting DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting DNS record: %s", resp.Msg)
	}

	return nil
}

func updateRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.UpdateRecord(
		ctx,
		dns.UpdateRecordRequest{
			Domain:   fmt.Sprintf("%v", d.Get("domain")),
			RecordID: d.Id(),
			Type:     fmt.Sprintf("%v", d.Get("type")),
			Name:     fmt.Sprintf("%v", d.Get("name")),
			Content:  fmt.Sprintf("%v", d.Get("content")),
			Priority: fmt.Sprintf("%v", d.Get("priority")),
		},
	)
	if err != nil {
		return diag.Errorf("Error updating DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error updating DNS record: %s", resp.Msg)
	}

	return nil
}

func importRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		if err := d.Set("domain", s[0]); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

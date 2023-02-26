// Package server provides the functions to create a Server resource via SiteHost API.
package dns

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/dns"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// ZoneResource returns a schema with the operations for Server resource.
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

// createResource is a function to create a new Server resource.
func createZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	domain := d.Get("name").(string)
	// The response don't have the domain name, so we need to get it from the request.
	_, err := client.CreateZone(ctx, dns.CreateZoneRequest{DomainName: domain})

	if err != nil {
		return diag.Errorf("Error creating domain: %s", err)
	}

	d.SetId(domain)

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	domain := d.Id()
	response, err := client.GetZone(ctx, dns.GetZoneRequest{DomainName: d.Id()})

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	domains := response.Return
	exist := false
	// iterate over the domains to find the one we are looking for.
	for _, d := range domains {
		if d.Name == domain {
			exist = true
			break
		}
	}

	if !exist {
		return diag.Errorf("Failed retrieving domain: %s", err)
	}

	d.SetId(domain)
	d.Set("name", domain)

	log.Printf("[INFO] Domain Name: %s", domain)

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	_, err := client.DeleteZone(ctx, dns.DeleteZoneRequest{DomainName: d.Id()})

	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	return nil
}

// Resource returns a schema with the operations for Server resource.
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

	domain := d.Get("domain").(string)
	name := d.Get("name").(string)
	domainRecord := helper.ConstructFqdn(name, domain)

	client := dns.New(conf.Client)
	_, err := client.AddRecord(
		ctx,
		dns.AddRecordRequest{
			Domain:   domain,
			Type:     d.Get("type").(string),
			Name:     domainRecord,
			Content:  d.Get("record").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
		},
	)

	if err != nil {
		return diag.Errorf("Error creating domain_record: %s", err)
	}

	//updateRecord(d, domainRecord)

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

	_, err := client.ListRecords(
		ctx,
		dns.ListRecordsRequest{Domain: d.Get("domain").(string)}, // TODO: This needs to be checked
	) // Check how handle the list

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	// updateRecordResource(d, domainRecord)

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	_, err := client.DeleteRecord(ctx, dns.DeleteRecordRequest{Domain: d.Get("domain").(string), RecordID: d.Id()})

	if err != nil {
		return diag.Errorf("Error deleting domain record: %s", err)
	}

	return nil
}

func updateRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	_, err := client.UpdateRecord(
		ctx,
		dns.UpdateRecordRequest{
			Domain:   d.Get("domain").(string),
			RecordID: d.Id(),
			Type:     d.Get("type").(string),
			Name:     d.Get("name").(string),
			Content:  d.Get("content").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
		},
	)

	if err != nil {
		return diag.Errorf("Error updating domain record: %s", err)
	}

	// updateRecord(d, domainRecord)

	return nil
}

func importRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}

// func updateRecord(d *schema.ResourceData, domainRecord *models.DNSRecord) {
// 	d.SetId(domainRecord.ID)
// 	d.Set("domain", domainRecord.Domain)
// 	d.Set("name", helper.DeconstructFqdn(domainRecord.Name, domainRecord.Domain))
// 	d.Set("record", domainRecord.Content)
// 	d.Set("type", domainRecord.Type)

// 	ttl, e := strconv.Atoi(domainRecord.TTL)
// 	if e == nil {
// 		d.Set("ttl", ttl)
// 	} else {
// 		d.Set("ttl", 3600)
// 	}

// 	priority, e := strconv.Atoi(domainRecord.Priority)
// 	if e == nil {
// 		d.Set("priority", priority)
// 	} else {
// 		d.Set("priority", 0)
// 	}

// 	d.Set("fqdn", helper.ConstructFqdn(domainRecord.Name, domainRecord.Domain))
// 	d.Set("change_date", domainRecord.ChangeDate)
// }

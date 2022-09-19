// Package server provides the functions to create a Server resource via SiteHost API.
package domain_record

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"log"
	"strconv"
	"strings"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new Server resource.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	domainRecord, err := client.DomainRecord.Add(
		ctx,
		&gosh.DomainRecord{
			ClientID: client.ClientID,
			Domain:   domain,
			Type:     d.Get("type").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
			Name:     constructFqdn(name, domain),
			Content:  d.Get("record").(string),
		},
	)

	if err != nil {
		return diag.Errorf("Error creating domain_record: %s", err)
	}
	if domainRecord == nil {
		return diag.Errorf("Failed creating domain_record: %s", domainRecord)
	}

	d.SetId(domainRecord.Id)
	d.Set("ttl", strconv.Itoa(d.Get("ttl").(int)))
	d.Set("fdqn", constructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("change_date", domainRecord.ChangeDate)
	d.Set("record", domainRecord.Content)

	log.Printf("[INFO] Domain Record: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	domainRecord, err := client.DomainRecord.Get(
		context.Background(),
		&gosh.Domain{
			ClientID: client.ClientID,
			Name:     d.Get("domain").(string),
		},
		d.Id(),
	)

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if domainRecord == nil {
		return diag.Errorf("Failed retrieving domain_record: %s", err)
	}

	d.SetId(domainRecord.Id)
	d.Set("domain", domainRecord.Domain)
	d.Set("name", deconstructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("record", domainRecord.Content)
	d.Set("type", domainRecord.Type)

	d.Set("ttl", domainRecord.TTL)
	d.Set("priority", domainRecord.Priority)

	d.Set("fdqn", constructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("change_date", domainRecord.ChangeDate)

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.CombinedConfig).Client

	_, err := client.DomainRecord.Delete(context.Background(), &gosh.DomainRecord{
		ClientID: client.ClientID,
		Name:     d.Get("name").(string),
		Domain:   d.Get("domain").(string),
		Id:       d.Id(),
	})

	if err != nil {
		return diag.Errorf("Error deleting domain: %s", err)
	}

	return nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func constructFqdn(name, domain string) string {
	if name == "@" {
		return domain
	}

	rn := strings.ToLower(name)
	domainSuffix := domain + "."
	if strings.HasSuffix(rn, domainSuffix) {
		rn = strings.TrimSuffix(rn, ".")
	} else {
		rn = strings.Join([]string{name, domain}, ".")
	}
	return rn
}

func deconstructFqdn(name, domain string) string {
	if name == domain {
		return "@"
	}

	rn := strings.ToLower(name)
	rn = strings.TrimSuffix(rn, ".")
	rn = strings.TrimSuffix(rn, "."+domain)

	return rn
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}

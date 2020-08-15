package sitehost

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/internal/api"
)

var resourceServer = schema.Resource{
	CreateContext: crud(resourceServerCreate),
	ReadContext:   crud(resourceServerRead),
	UpdateContext: crud(resourceServerUpdate),
	DeleteContext: crud(resourceServerDelete),
	Schema: map[string]*schema.Schema{
		nameKey: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		}, passwordKey: {
			Type:      schema.TypeString,
			Sensitive: true,
			Computed:  true,
		}, ipsKey: {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		}, labelKey: {
			Type:     schema.TypeString,
			Required: true,
		}, locationKey: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		}, productCodeKey: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		}, imageKey: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	},
}

func unmarshalIPs(s api.Server) []string {
	ips := make([]string, 0, len(s.IPs))
	for _, ip := range s.IPs {
		ips = append(ips, ip.IPAddr.String())
	}
	return ips
}

func marshalIPs(v interface{}, ok bool) []api.ServerProvisionOption {
	if !ok {
		return []api.ServerProvisionOption{api.ParamIP{}}
	}
	ips := v.([]interface{})
	if len(ips) == 0 {
		return []api.ServerProvisionOption{api.ParamIP{}}
	}
	opts := make([]api.ServerProvisionOption, 0, len(ips))
	for _, ip := range ips {
		opts = append(opts, api.ParamIP(net.ParseIP(ip.(string))))
	}
	return opts
}

func resourceServerCreate(ctx context.Context, c api.Client, d *schema.ResourceData) error {
	var opts []api.ServerProvisionOption
	opts = append(opts, marshalIPs(d.GetOk(ipsKey))...)
	if v, ok := d.GetOk(nameKey); ok {
		opts = append(opts, api.ParamName(v.(string)))
	}
	id, job, name, password, addr, err := c.Server().Provision(ctx,
		d.Get(labelKey).(string),
		d.Get(locationKey).(string),
		d.Get(productCodeKey).(string),
		d.Get(imageKey).(string), opts...)
	if err != nil {
		return err
	}
	if err = wait(ctx, c.Job(), job); err != nil {
		return err
	}
	d.SetId(utoa(id))
	ips := make([]string, 0, len(addr))
	for _, ip := range addr {
		ips = append(ips, ip.String())
	}
	for k, v := range map[string]interface{}{
		nameKey:     name,
		passwordKey: password,
		ipsKey:      ips,
	} {
		if err = d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func resourceServerRead(ctx context.Context, c api.Client, d *schema.ResourceData) error {
	s, err := c.Server().Get(ctx, d.Get(nameKey).(string))
	switch err {
	case api.ErrNotFound:
		d.SetId("")
		return nil
	case nil:
	default:
		return err
	}
	d.SetId(utoa(s.ID))
	for k, v := range map[string]interface{}{
		ipsKey:         unmarshalIPs(s),
		labelKey:       s.Label,
		locationKey:    s.Location,
		productCodeKey: s.ProductCode,
	} {
		if err = d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func resourceServerUpdate(ctx context.Context, c api.Client, d *schema.ResourceData) error {
	panic(nil)
}

func resourceServerDelete(ctx context.Context, c api.Client, d *schema.ResourceData) error {
	job, err := c.Server().Delete(ctx, d.Get(nameKey).(string))
	if err != nil {
		return err
	}
	return wait(ctx, c.Job(), job)
}

// wait blocks until the job completes.
func wait(ctx context.Context, c *api.JobClient, job uint) error {
	for {
		j, err := c.Get(ctx, job, api.DaemonJob)
		switch {
		case err != nil:
			return err
		case j.State == api.CompleteJob:
			return nil
		}
	}
}

package sitehost

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/terraform-provider-sitehost/internal/api"
)

func crud(f func(ctx context.Context, c api.Client, d *schema.ResourceData) error) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		if err := f(ctx, m.(api.Client), d); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

// utoa is equivalent to strconv.FormatUint(uint64(u), 10).
//
// It is based on strconv.Itoa.
func utoa(u uint) string { return strconv.FormatUint(uint64(u), 10) }

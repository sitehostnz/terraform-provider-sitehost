package security_groups

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server/securitygroups"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for security groups.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}


// updateResource is a function to update a security group.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	// Only process updates if label or rules have changed
	if !d.HasChange("label") && !d.HasChange("rules") {
		return nil
	}

	// Get the rules from the schema
	rulesRaw, ok := d.Get("rules").(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to get rules from schema")
	}

	// Process inbound rules
	var rulesIn []securitygroups.RuleIn
	if inRulesRaw, ok := rulesRaw["in"].([]interface{}); ok {
		for _, ruleRaw := range inRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			rulesIn = append(rulesIn, securitygroups.RuleIn{
				Enabled:         rule["enabled"].(bool),
				Action:          rule["action"].(string),
				Protocol:        rule["protocol"].(string),
				SourceIP:        rule["src_ip"].(string),
				DestinationPort: rule["dest_port"].(string),
			})
		}
	}

	// Process outbound rules
	var rulesOut []securitygroups.RuleOut
	if outRulesRaw, ok := rulesRaw["out"].([]interface{}); ok {
		for _, ruleRaw := range outRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			rulesOut = append(rulesOut, securitygroups.RuleOut{
				Enabled:         rule["enabled"].(bool),
				Action:          rule["action"].(string),
				Protocol:        rule["protocol"].(string),
				DestinationIP:   rule["dest_ip"].(string),
				DestinationPort: rule["dest_port"].(string),
			})
		}
	}

	opts := securitygroups.UpdateRequest{
		Name: d.Id(),
		Params: securitygroups.ParamsOptions{
			Label:    fmt.Sprint(d.Get("label")),
			RulesIn:  rulesIn,
			RulesOut: rulesOut,
		},
	}

	res, err := client.Update(ctx, opts)
	if err != nil {
		return diag.Errorf("Error updating security group: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error updating security group: %s", res.Msg)
	}

	return nil
}

// deleteResource is a function to delete a security group.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	resp, err := client.Delete(context.Background(), securitygroups.DeleteRequest{
		Name: d.Id(),
	})

	if err != nil {
		return diag.Errorf("Error deleting security group: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting security group: %s", resp.Msg)
	}

	return nil
}

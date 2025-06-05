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
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new security group.
func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	opts := securitygroups.AddRequest{
		Label: fmt.Sprint(d.Get("label")),
	}

	res, err := client.Add(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating security group: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error creating security group: %s", res.Msg)
	}

	// Set data
	d.SetId(res.Return.Name)
	if err := d.Set("name", res.Return.Name); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Security Group Name: %s", d.Id())

	return nil
}

// updateResource is a function to update a security group.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)

	// Only process updates if label or rules have changed
	if !d.HasChange("label") && !d.HasChange("rules_in") && !d.HasChange("rules_out") {
		return nil
	}

	// Process inbound rules
	var rulesIn []securitygroups.RuleIn
	if inRulesRaw, ok := d.Get("rules_in").([]interface{}); ok {
		for _, ruleRaw := range inRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			enabled := true
			if e, ok := rule["enabled"].(bool); ok {
				enabled = e
			}
			rulesIn = append(rulesIn, securitygroups.RuleIn{
				Enabled:         enabled,
				Action:          rule["action"].(string),
				Protocol:        rule["protocol"].(string),
				SourceIP:        rule["src_ip"].(string),
				DestinationPort: rule["dest_port"].(string),
			})
		}
	}

	// Process outbound rules
	var rulesOut []securitygroups.RuleOut
	if outRulesRaw, ok := d.Get("rules_out").([]interface{}); ok {
		for _, ruleRaw := range outRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			enabled := true
			if e, ok := rule["enabled"].(bool); ok {
				enabled = e
			}
			rulesOut = append(rulesOut, securitygroups.RuleOut{
				Enabled:         enabled,
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

	return readResource(ctx, d, meta)
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

// readResource is a function to read a security group.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := securitygroups.New(conf.Client)
	resp, err := client.Get(ctx, securitygroups.GetRequest{
		Name: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error reading security group: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error reading security group: %s", resp.Msg)
	}

	if err := d.Set("label", resp.Return.Label); err != nil {
		return diag.FromErr(err)
	}

	// Set inbound rules
	rulesIn := make([]map[string]interface{}, len(resp.Return.Rules.In))
	for i, rule := range resp.Return.Rules.In {
		rulesIn[i] = map[string]interface{}{
			"enabled":   rule.Enabled,
			"action":    rule.Action,
			"protocol":  rule.Protocol,
			"src_ip":    rule.SrcIP,
			"dest_port": fmt.Sprint(rule.DestPort),
		}
	}
	if err := d.Set("rules_in", rulesIn); err != nil {
		return diag.FromErr(err)
	}

	// Set outbound rules
	rulesOut := make([]map[string]interface{}, len(resp.Return.Rules.Out))
	for i, rule := range resp.Return.Rules.Out {
		rulesOut[i] = map[string]interface{}{
			"enabled":   rule.Enabled,
			"action":    rule.Action,
			"protocol":  rule.Protocol,
			"dest_ip":   rule.DestIP,
			"dest_port": fmt.Sprint(rule.DestPort),
		}
	}
	if err := d.Set("rules_out", rulesOut); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

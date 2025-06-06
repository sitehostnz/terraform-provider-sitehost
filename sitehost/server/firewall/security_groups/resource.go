package security_groups

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server/firewall/security_groups"
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

	client := security_groups.New(conf.Client)

	opts := security_groups.AddRequest{
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

	if d.Get("rules_in") != nil || d.Get("rules_out") != nil {
		if err := updateResource(ctx, d, meta); err != nil {
			return err
		}
	}

	return nil
}

// updateResource is a function to update a security group.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := security_groups.New(conf.Client)

	// Process inbound rules
	var rulesIn []security_groups.RuleIn
	if inRulesRaw, ok := d.Get("rules_in").([]interface{}); ok {
		for _, ruleRaw := range inRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			enabled := true
			if e, ok := rule["enabled"].(bool); ok {
				enabled = e
			}
			rulesIn = append(rulesIn, security_groups.RuleIn{
				Enabled:         enabled,
				Action:          rule["action"].(string),
				Protocol:        rule["protocol"].(string),
				SourceIP:        rule["src_ip"].(string),
				DestinationPort: rule["dest_port"].(string),
			})
		}
	}

	// Process outbound rules
	var rulesOut []security_groups.RuleOut
	if outRulesRaw, ok := d.Get("rules_out").([]interface{}); ok {
		for _, ruleRaw := range outRulesRaw {
			rule := ruleRaw.(map[string]interface{})
			enabled := true
			if e, ok := rule["enabled"].(bool); ok {
				enabled = e
			}
			rulesOut = append(rulesOut, security_groups.RuleOut{
				Enabled:         enabled,
				Action:          rule["action"].(string),
				Protocol:        rule["protocol"].(string),
				DestinationIP:   rule["dest_ip"].(string),
				DestinationPort: rule["dest_port"].(string),
			})
		}
	}

	opts := security_groups.UpdateRequest{
		Name: d.Id(),
		Params: security_groups.ParamsOptions{
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

	client := security_groups.New(conf.Client)

	resp, err := client.Delete(context.Background(), security_groups.DeleteRequest{
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

	client := security_groups.New(conf.Client)
	resp, err := client.Get(ctx, security_groups.GetRequest{
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

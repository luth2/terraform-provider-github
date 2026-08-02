package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v89/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGithubEnterpriseBillingCostCenter() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a GitHub Enterprise billing cost center.",
		CreateContext: resourceGithubEnterpriseBillingCostCenterCreate,
		ReadContext:   resourceGithubEnterpriseBillingCostCenterRead,
		UpdateContext: resourceGithubEnterpriseBillingCostCenterUpdate,
		DeleteContext: resourceGithubEnterpriseBillingCostCenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGithubEnterpriseBillingCostCenterImport,
		},

		Schema: map[string]*schema.Schema{
			"enterprise_slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The slug of the enterprise.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cost center.",
			},
			"cost_center_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cost center assigned by GitHub.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the cost center (`active` or `deleted`).",
			},
			"azure_subscription": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Azure subscription linked to this cost center, if any.",
			},
			"resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resources assigned to this cost center.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type (e.g. `organization`, `repository`, `user`).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the resource.",
						},
					},
				},
			},
		},
	}
}

func resourceGithubEnterpriseBillingCostCenterCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterpriseSlug := d.Get("enterprise_slug").(string)
	req := github.CostCenterRequest{
		Name: d.Get("name").(string),
	}

	cc, _, err := client.Enterprise.CreateCostCenter(ctx, enterpriseSlug, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", enterpriseSlug, cc.ID))
	if err := d.Set("cost_center_id", cc.ID); err != nil {
		return diag.FromErr(err)
	}

	return resourceGithubEnterpriseBillingCostCenterRead(ctx, d, meta)
}

func resourceGithubEnterpriseBillingCostCenterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterpriseSlug, costCenterID, err := parseBillingCostCenterID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cc, resp, err := client.Enterprise.GetCostCenter(ctx, enterpriseSlug, costCenterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("enterprise_slug", enterpriseSlug); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cost_center_id", cc.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", cc.Name); err != nil {
		return diag.FromErr(err)
	}
	if cc.State != nil {
		if err := d.Set("state", *cc.State); err != nil {
			return diag.FromErr(err)
		}
	}
	if cc.AzureSubscription != nil {
		if err := d.Set("azure_subscription", *cc.AzureSubscription); err != nil {
			return diag.FromErr(err)
		}
	}

	var resources []map[string]string
	for _, r := range cc.Resources {
		resources = append(resources, map[string]string{
			"type": r.Type,
			"name": r.Name,
		})
	}
	if err := d.Set("resources", resources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGithubEnterpriseBillingCostCenterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterpriseSlug, costCenterID, err := parseBillingCostCenterID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := github.CostCenterRequest{
		Name: d.Get("name").(string),
	}

	_, _, err = client.Enterprise.UpdateCostCenter(ctx, enterpriseSlug, costCenterID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGithubEnterpriseBillingCostCenterRead(ctx, d, meta)
}

func resourceGithubEnterpriseBillingCostCenterDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterpriseSlug, costCenterID, err := parseBillingCostCenterID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client.Enterprise.DeleteCostCenter(ctx, enterpriseSlug, costCenterID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGithubEnterpriseBillingCostCenterImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	_, _, err := parseBillingCostCenterID(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid import ID %q: expected <enterprise_slug>/<cost_center_id>", d.Id())
	}

	diags := resourceGithubEnterpriseBillingCostCenterRead(ctx, d, meta)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading cost center: %s", diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

func parseBillingCostCenterID(id string) (enterpriseSlug, costCenterID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid cost center ID %q: expected <enterprise_slug>/<cost_center_id>", id)
	}
	return parts[0], parts[1], nil
}

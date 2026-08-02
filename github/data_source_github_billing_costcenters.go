package github

import (
	"context"

	"github.com/google/go-github/v89/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGithubBillingCostCenters() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all cost centers for a GitHub Enterprise.",
		ReadContext: dataSourceGithubBillingCostCentersRead,

		Schema: map[string]*schema.Schema{
			"enterprise_slug": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the enterprise.",
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "deleted"}, false),
				Description:  "Filter cost centers by state. Valid values: `active`, `deleted`.",
			},
			"cost_centers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_subscription": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGithubBillingCostCentersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client

	enterpriseSlug := d.Get("enterprise_slug").(string)

	opts := &github.ListCostCenterOptions{}
	if v, ok := d.GetOk("state"); ok {
		s := v.(string)
		opts.State = &s
	}

	costCenters, _, err := client.Enterprise.ListCostCenters(ctx, enterpriseSlug, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(enterpriseSlug)

	var result []map[string]any
	for _, cc := range costCenters.CostCenters {
		var resources []map[string]string
		for _, r := range cc.Resources {
			resources = append(resources, map[string]string{
				"type": r.Type,
				"name": r.Name,
			})
		}

		item := map[string]any{
			"id":        cc.ID,
			"name":      cc.Name,
			"resources": resources,
		}
		if cc.State != nil {
			item["state"] = *cc.State
		}
		if cc.AzureSubscription != nil {
			item["azure_subscription"] = *cc.AzureSubscription
		}
		result = append(result, item)
	}

	if err := d.Set("cost_centers", result); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

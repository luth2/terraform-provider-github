package github

import (
	"context"

	"github.com/google/go-github/v89/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGithubBillingOrganizationUsage() *schema.Resource {
	return &schema.Resource{
		Description: "Gets a billing usage report for a GitHub organization (requires enhanced billing platform).",
		ReadContext: dataSourceGithubBillingOrganizationUsageRead,

		Schema: map[string]*schema.Schema{
			"year": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2000, 9999),
				Description:  "Filter results to a specific year (e.g. 2025). Defaults to the current year.",
			},
			"month": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 12),
				Description:  "Filter results to a specific month (1–12).",
			},
			"day": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 31),
				Description:  "Filter results to a specific day (1–31).",
			},
			"hour": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 23),
				Description:  "Filter results to a specific hour (0–23).",
			},
			"usage_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quantity": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"unit_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price_per_unit": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"gross_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"discount_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"net_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGithubBillingOrganizationUsageRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	org := meta.(*Owner).name

	opts := &github.UsageReportOptions{}
	if v, ok := d.GetOk("year"); ok {
		year := v.(int)
		opts.Year = &year
	}
	if v, ok := d.GetOk("month"); ok {
		month := v.(int)
		opts.Month = &month
	}
	if v, ok := d.GetOk("day"); ok {
		day := v.(int)
		opts.Day = &day
	}
	if v, ok := d.GetOk("hour"); ok {
		hour := v.(int)
		opts.Hour = &hour
	}

	report, _, err := client.Billing.GetOrganizationUsageReport(ctx, org, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(org)

	var items []map[string]any
	for _, item := range report.UsageItems {
		m := map[string]any{
			"date":            item.Date,
			"product":         item.Product,
			"sku":             item.SKU,
			"quantity":        item.Quantity,
			"unit_type":       item.UnitType,
			"price_per_unit":  item.PricePerUnit,
			"gross_amount":    item.GrossAmount,
			"discount_amount": item.DiscountAmount,
			"net_amount":      item.NetAmount,
		}
		if item.RepositoryName != nil {
			m["repository_name"] = *item.RepositoryName
		}
		if item.OrganizationName != nil {
			m["organization_name"] = *item.OrganizationName
		}
		items = append(items, m)
	}

	if err := d.Set("usage_items", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

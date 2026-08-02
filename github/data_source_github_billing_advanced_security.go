package github

import (
	"context"

	"github.com/google/go-github/v89/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGithubBillingAdvancedSecurity() *schema.Resource {
	return &schema.Resource{
		Description: "Gets GitHub Advanced Security active committer counts for an organization.",
		ReadContext: dataSourceGithubBillingAdvancedSecurityRead,

		Schema: map[string]*schema.Schema{
			"advanced_security_product": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"code_security", "secret_protection"}, false),
				Description:  "The security product to filter committers for. Valid values: `code_security`, `secret_protection`.",
			},
			"total_advanced_security_committers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"maximum_advanced_security_committers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"purchased_advanced_security_committers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_security_committers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"advanced_security_committers_breakdown": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_login": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_pushed_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_pushed_email": {
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

func dataSourceGithubBillingAdvancedSecurityRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Owner).v3client
	org := meta.(*Owner).name

	opts := &github.ActiveCommittersListOptions{
		ListOptions: github.ListOptions{
			PerPage: meta.(*Owner).maxPerPage,
		},
	}
	if v, ok := d.GetOk("advanced_security_product"); ok {
		p := v.(string)
		opts.AdvancedSecurityProduct = &p
	}

	var allRepos []*github.RepositoryActiveCommitters
	var totalCommitters, totalCount, maxCommitters, purchasedCommitters *int

	for {
		result, resp, err := client.Billing.GetOrganizationAdvancedSecurityActiveCommitters(ctx, org, opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if totalCommitters == nil {
			totalCommitters = result.TotalAdvancedSecurityCommitters
			totalCount = result.TotalCount
			maxCommitters = result.MaximumAdvancedSecurityCommitters
			purchasedCommitters = result.PurchasedAdvancedSecurityCommitters
		}

		allRepos = append(allRepos, result.Repositories...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	d.SetId(org)

	if totalCommitters != nil {
		if err := d.Set("total_advanced_security_committers", *totalCommitters); err != nil {
			return diag.FromErr(err)
		}
	}
	if totalCount != nil {
		if err := d.Set("total_count", *totalCount); err != nil {
			return diag.FromErr(err)
		}
	}
	if maxCommitters != nil {
		if err := d.Set("maximum_advanced_security_committers", *maxCommitters); err != nil {
			return diag.FromErr(err)
		}
	}
	if purchasedCommitters != nil {
		if err := d.Set("purchased_advanced_security_committers", *purchasedCommitters); err != nil {
			return diag.FromErr(err)
		}
	}

	var repos []map[string]any
	for _, r := range allRepos {
		var breakdown []map[string]string
		for _, b := range r.AdvancedSecurityCommittersBreakdown {
			breakdown = append(breakdown, map[string]string{
				"user_login":        b.UserLogin,
				"last_pushed_date":  b.LastPushedDate,
				"last_pushed_email": b.LastPushedEmail,
			})
		}
		repos = append(repos, map[string]any{
			"name":                                   r.Name,
			"advanced_security_committers":           r.AdvancedSecurityCommitters,
			"advanced_security_committers_breakdown": breakdown,
		})
	}

	if err := d.Set("repositories", repos); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

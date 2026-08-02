---
page_title: "github_billing_advanced_security (Data Source) - GitHub"
description: |-
  Gets GitHub Advanced Security active committer counts for an organization.
---

# github_billing_advanced_security (Data Source)

Use this data source to retrieve GitHub Advanced Security (GHAS) active committer counts per repository for the configured organization.

## Example Usage

```terraform
data "github_billing_advanced_security" "all" {}

data "github_billing_advanced_security" "code_security" {
  advanced_security_product = "code_security"
}
```

## Argument Reference

- `advanced_security_product` - (Optional) Filter committers for a specific security product. Valid values: `code_security`, `secret_protection`.

## Attributes Reference

- `total_advanced_security_committers` - Total number of active GHAS committers across all repositories.
- `total_count` - Total number of unique committers.
- `maximum_advanced_security_committers` - Maximum number of GHAS committers allowed under the license.
- `purchased_advanced_security_committers` - Number of GHAS seats purchased.
- `repositories` - Per-repository committer breakdown.
  - `name` - The repository name.
  - `advanced_security_committers` - Number of active GHAS committers in this repository.
  - `advanced_security_committers_breakdown` - Per-user breakdown.
    - `user_login` - The GitHub username.
    - `last_pushed_date` - The date of the user's last push.
    - `last_pushed_email` - The email used for the user's last push.

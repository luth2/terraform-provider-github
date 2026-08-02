---
page_title: "github_billing_organization_usage (Data Source) - GitHub"
description: |-
  Gets the billing usage report for a GitHub organization.
---

# github_billing_organization_usage (Data Source)

Use this data source to retrieve the billing usage report for the configured GitHub organization. Requires the enhanced billing platform.

## Example Usage

```terraform
data "github_billing_organization_usage" "current" {}

data "github_billing_organization_usage" "monthly" {
  year  = 2025
  month = 3
}
```

## Argument Reference

- `year` - (Optional) Filter results to a specific year (e.g. `2025`). Defaults to the current year.
- `month` - (Optional) Filter results to a specific month (1–12).
- `day` - (Optional) Filter results to a specific day (1–31).
- `hour` - (Optional) Filter results to a specific hour (0–23).

## Attributes Reference

- `usage_items` - List of usage line items.
  - `date` - The date of the usage.
  - `product` - The product name (e.g. `Actions`, `Packages`).
  - `sku` - The SKU identifier.
  - `quantity` - The quantity consumed.
  - `unit_type` - The unit of measurement.
  - `price_per_unit` - The price per unit.
  - `gross_amount` - The gross amount charged.
  - `discount_amount` - The discount applied.
  - `net_amount` - The net amount after discount.
  - `repository_name` - The repository associated with the usage, if applicable.
  - `organization_name` - The organization associated with the usage, if applicable.

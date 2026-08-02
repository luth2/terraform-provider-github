---
page_title: "github_billing_cost_centers (Data Source) - GitHub"
description: |-
  Lists all billing cost centers for a GitHub Enterprise.
---

# github_billing_cost_centers (Data Source)

Use this data source to list all billing cost centers for a GitHub Enterprise (requires the enhanced billing platform).

## Example Usage

```terraform
data "github_billing_cost_centers" "all" {
  enterprise_slug = "my-enterprise"
}

data "github_billing_cost_centers" "active" {
  enterprise_slug = "my-enterprise"
  state           = "active"
}
```

## Argument Reference

- `enterprise_slug` - (Required) The slug of the enterprise.
- `state` - (Optional) Filter cost centers by state. Valid values: `active`, `deleted`.

## Attributes Reference

- `cost_centers` - List of cost centers.
  - `id` - The ID of the cost center.
  - `name` - The name of the cost center.
  - `state` - The state of the cost center (`active` or `deleted`).
  - `azure_subscription` - The linked Azure subscription, if any.
  - `resources` - List of resources assigned to the cost center.
    - `type` - The resource type (e.g. `organization`, `repository`, `user`).
    - `name` - The name of the resource.

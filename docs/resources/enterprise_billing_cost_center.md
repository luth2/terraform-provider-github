---
page_title: "github_enterprise_billing_cost_center (Resource) - GitHub"
description: |-
  Manages a GitHub Enterprise billing cost center.
---

# github_enterprise_billing_cost_center (Resource)

Manage billing cost centers for a GitHub Enterprise. Requires the enhanced billing platform and an account with enterprise admin permissions.

## Example Usage

```terraform
resource "github_enterprise_billing_cost_center" "engineering" {
  enterprise_slug = "my-enterprise"
  name            = "Engineering"
}
```

## Argument Reference

- `enterprise_slug` - (Required, Forces new resource) The slug of the enterprise.
- `name` - (Required) The name of the cost center.

## Attributes Reference

- `cost_center_id` - The ID of the cost center assigned by GitHub.
- `state` - The state of the cost center (`active` or `deleted`).
- `azure_subscription` - The linked Azure subscription, if any.
- `resources` - List of resources assigned to this cost center.
  - `type` - The resource type (e.g. `organization`, `repository`, `user`).
  - `name` - The name of the resource.

## Import

Cost centers can be imported using `<enterprise_slug>/<cost_center_id>`:

```shell
terraform import github_enterprise_billing_cost_center.engineering my-enterprise/cc-abc123
```

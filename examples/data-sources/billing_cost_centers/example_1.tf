data "github_billing_cost_centers" "all" {
  enterprise_slug = "my-enterprise"
}

data "github_billing_cost_centers" "active" {
  enterprise_slug = "my-enterprise"
  state           = "active"
}

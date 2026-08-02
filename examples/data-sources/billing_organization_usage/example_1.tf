data "github_billing_organization_usage" "current" {}

data "github_billing_organization_usage" "march" {
  year  = 2025
  month = 3
}

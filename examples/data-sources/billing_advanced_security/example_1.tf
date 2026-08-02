data "github_billing_advanced_security" "all" {}

data "github_billing_advanced_security" "code_security" {
  advanced_security_product = "code_security"
}

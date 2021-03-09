

locals {
  VAULT_ADDR = "http://88.97.2.109:8200"
}

// Vault auth during Planning time - can only have static query elements
data "external" "vault_login_plan" {
  program = ["bash", "${path.module}/vault_login.sh"]
  query = {
    role             = "workspace_role"
    VAULT_ADDR       = local.VAULT_ADDR
    VAULT_LOGIN_PATH = "v1/auth/tfe-auth/login"
  }
}

// Vault auth during Apply time - must have a dynamic element
data "external" "vault_login_apply" {
  program = ["bash", "${path.module}/vault_login.sh"]
  query = {
    role             = "workspace_role"
    VAULT_ADDR       = local.VAULT_ADDR
    VAULT_LOGIN_PATH = "v1/auth/tfe-auth/login"
    always_run       = timestamp()
  }
}

provider "vault" {
  address    = local.VAULT_ADDR
  token      = data.external.vault_login_apply == null ? data.external.vault_login_plan.result.VAULT_TOKEN : data.external.vault_login_apply.result.VAULT_TOKEN
  token_name = "terraform-${var.TFE_RUN_ID}"
}

# the below code block does not work and is here as an example for a provider improvement suggestion.
# provider "vault" {
#   address    = "http://88.97.2.109:8200"
#   token_name = "terraform-${var.TFE_RUN_ID}"
#   auth_login {
#     path = "auth/tfe-auth/login"
#     parameters = {
#       role      = "workspace_role"
#       workspace = var.TFC_WORKSPACE_NAME
#       run-id    = var.TFE_RUN_ID
#       atlas-token = "env(ATLAS_TOKEN)"
#     }
#   }
# }

// just a test.
resource "vault_generic_secret" "example" {
  path = "secret/hello"

  data_json = <<EOT
{
  "foo":   "bar",
  "pizza": "crackers"
}
EOT
}

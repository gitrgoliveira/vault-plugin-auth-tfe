

// Vault auth during Planning time - can only have static query elements
data "external" "vault_login" {
  program = ["bash", "${path.module}/vault_login.sh"]
  query = {
    role = "workspace_role"
    VAULT_ADDR = "http://88.97.2.109:8200"
  }
}

// Vault auth during Apply time
# data "external" "vault_login2" {
#   program = ["bash", "${path.module}/vault_login.sh"]
#   query = {
#     always_run = timestamp()
#     role = "workspace_role"
#   }
#   depends_on = [ data.external.vault_login ]
# }


provider "vault" {
  address    = "http://88.97.2.109:8200"
  token      = data.external.vault_login.result.VAULT_TOKEN
  token_name = "terraform-${var.TFE_RUN_ID}"
}

// The atlas token keeps changing as cannot be updated directly in the provider.
# data "external" "atlas_token" {
#   program = ["bash", "${path.module}/atlas_token.sh"]
# }
# data "external" "atlas_token2" {
#   program = ["bash", "${path.module}/atlas_token.sh"]
#   query = {
#     always_run = timestamp()
#   }
# }
# provider "vault" {
#   address    = "http://88.97.2.109:8200"
#   token_name = "terraform-${var.TFE_RUN_ID}"
#   auth_login {
#     path = "auth/tfe-auth/login"
#     parameters = {
#       role      = "workspace_role"
#       workspace = var.TFC_WORKSPACE_NAME
#       run-id    = var.TFE_RUN_ID
#       atlas-token = data.external.atlas_token.result.ATLAS_TOKEN
#     }
#   }
# }

resource "vault_generic_secret" "example" {
  path  = "secret/hello"

  data_json = <<EOT
{
  "foo":   "bar",
  "pizza": "crackers"
}
EOT
}

# data "vault_generic_secret" "hello_secret" {
#   path = "secret/hello"
# }

# resource "vault_namespace" "vault-team" {
#   path     = "vault-team"
# }
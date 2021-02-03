

// Vault auth during Planning time - can only have static query elements
data "external" "vault_login" {
  program = ["bash", "${path.module}/vault_login.sh"]
  query = {
    role = "workspace_role"
    VAULT_ADDR = "http://88.97.2.109:8200"
  }
}

provider "vault" {
  address    = "http://88.97.2.109:8200"
  token      = data.external.vault_login.result.VAULT_TOKEN
  token_name = "terraform-${var.TFE_RUN_ID}"
}

resource "vault_generic_secret" "example" {
  path  = "secret/hello"

  data_json = <<EOT
{
  "foo":   "bar",
  "pizza": "crackers"
}
EOT
}

#  helpers to figure out the TFC Agent file path
provider "vault" {
  address    = var.VAULT_ADDR
  token_name = "terraform-${var.TFE_RUN_ID}"
  // auth_login {
  //   path = "auth/tfe-auth/login"
  //   parameters = {
  //     role      = "workspace_role"
  //     workspace = var.TFC_WORKSPACE_NAME
  //     run-id    = var.TFE_RUN_ID
  //     # For code that is running within TFC/TFE or using an external agent
  //     tfe-credentials-file = try(filebase64("${path.cwd}/../.terraformrc"),
  //                                 filebase64("/tmp/cli.tfrc"))
  //   }
  // }
}

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

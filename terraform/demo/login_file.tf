#  helpers to figure out the TFC Agent file path
locals {
  file_path_apply = "/root/.tfc-agent/component/terraform/runs/${var.TFE_RUN_ID}.apply/cli.tfrc"
  file_path_plan  = "/root/.tfc-agent/component/terraform/runs/${var.TFE_RUN_ID}.plan/cli.tfrc"
  # file_path       = fileexists(local.file_path_apply) ? local.file_path_apply : local.file_path_plan
}

provider "vault" {
  address    = "http://88.97.2.109:8200"
  token_name = "terraform-${var.TFE_RUN_ID}"
  auth_login {
    path = "auth/tfe-auth/login"
    parameters = {
      role      = "workspace_role"
      workspace = var.TFC_WORKSPACE_NAME
      run-id    = var.TFE_RUN_ID
      # For code that is running within TFC/TFE
      tfe-credentials-file = filebase64("/tmp/cli.tfrc")

      # For code that is using a TFC/Runner
      # tfe-credentials-file = filebase64(local.file_path)

      # For code to run in both
      # tfe-credentials-file = fileexists("/tmp/cli.tfrc") ? filebase64("/tmp/cli.tfrc") : filebase64(local.file_path)
    }
  }
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

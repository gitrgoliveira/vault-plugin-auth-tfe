terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.0.0"
    }
  }
}

locals {
  workspace = split("/",var.TFC_WORKSPACE_SLUG)[1]
  organization = split("/",var.TFC_WORKSPACE_SLUG)[0]
}

provider "null" {
  # Configuration options
}
resource "null_resource" "run_info" {
  provisioner "local-exec" {
    command = "curl --header \"Authorization: Bearer $ATLAS_TOKEN\" --header \"Content-Type: application/vnd.api+json\"  $TF_VAR_ATLAS_ADDRESS/api/v2/runs/$ATLAS_RUN_ID"
  }
  triggers = {
    always_run = timestamp()
  }

}

resource "null_resource" "workspace_info" {
  provisioner "local-exec" {
    command = "curl --header \"Authorization: Bearer $ATLAS_TOKEN\" --header \"Content-Type: application/vnd.api+json\"  $TF_VAR_ATLAS_ADDRESS/api/v2/organizations/${local.organization}/workspaces/${local.workspace}"
  }
  triggers = {
    always_run = timestamp()
  }

}
resource "null_resource" "account_details" {
  provisioner "local-exec" {
    command = "curl --header \"Authorization: Bearer $ATLAS_TOKEN\" --header \"Content-Type: application/vnd.api+json\" --request GET $TF_VAR_ATLAS_ADDRESS/api/v2/account/details"
  }
  triggers = {
    always_run = timestamp()
  }
}


resource "null_resource" "sleep" {
  provisioner "local-exec" {
    command = "env && sleep 300"
  }
  triggers = {
    always_run = timestamp()
  }
  depends_on = [null_resource.account_details,
  null_resource.run_info,
  null_resource.workspace_info]

}


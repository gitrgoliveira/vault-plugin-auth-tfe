# terraform {
#   backend "remote" {
#     hostname     = "tfe.ric.gcp.hashidemos.io"
#     organization = "org2"
#     workspaces {
#       name = "vault-login-demo"
#     }
#   }
# }

terraform {
  backend "remote" {
    organization = "hc-emea-sentinel-demo"
    workspaces {
      name = "vault-login-demo"
    }
  }

  required_providers {
    vault = {
      source = "hashicorp/vault"
    }

  }
}
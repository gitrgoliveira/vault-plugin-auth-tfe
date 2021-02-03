terraform {
  backend "remote" {
    hostname     = "tfe.ric.gcp.hashidemos.io"
    organization = "org2"
    workspaces {
      name = "tfe-gcp-test-network"
    }
  }
}
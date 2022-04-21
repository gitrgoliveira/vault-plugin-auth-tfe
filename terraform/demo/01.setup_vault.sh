#! /bin/bash
source helper.sh

vault policy write terraform-policy - << EOF
path "auth/token/create" {
    capabilities = ["update"]
}
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

path "secret/data/*" {
  capabilities = ["read","create", "delete", "update"]
}
path "secret/*" {
    capabilities = ["read", "create", "delete", "update"]
}

path "aws/sts/deploy" {
  capabilities = ["read"]
}

EOF

vault auth enable -path=tfe-auth vault-plugin-auth-tfe
vault write auth/tfe-auth/config organization=hc-emea-sentinel-demo use_run_status=true
# vault write auth/tfe-auth/config organization=org2 \
#     terraform_host=https://tfe.ric.gcp.hashidemos.io

vault read auth/tfe-auth/config
vault write auth/tfe-auth/role/workspace_role workspaces=* \
    policies=default,terraform-policy

vault read auth/tfe-auth/role/workspace_role

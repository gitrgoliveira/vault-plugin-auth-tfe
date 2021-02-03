
export VAULT_ADDR='http://88.97.2.109:8200'
vault policy write terraform-policy - << EOF
path "auth/token/create" {
    capabilities = ["update"]
}
path "secret/data/*" {
  capabilities = ["read","create", "update"]
}

path "secret/*" {
    capabilities = ["read", "create", "update"]
}
EOF
vault kv put secret/hello foo=world


vault auth enable -path=tfe-auth vault-plugin-auth-tfe
# vault write auth/tfe-auth/config organization=hc-emea-sentinel-demo
vault write auth/tfe-auth/config organization=org2 \
    terraform_host=tfe.ric.gcp.hashidemos.io

vault read auth/tfe-auth/config
vault write auth/tfe-auth/role/workspace_role workspaces=* \
    policies=default,terraform-policy

vault read auth/tfe-auth/role/workspace_role

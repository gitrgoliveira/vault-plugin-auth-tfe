#! /bin/bash
# Exit if any of the intermediate steps fail
set -e
# jq is not present by default in TFC/TFE
curl -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod a+x jq

VAULT_PATH="v1/$VAULT_LOGIN_PATH"

export VAULT_TOKEN=$(curl -X PUT -H "X-Vault-Request: true" \
  -H "X-Vault-Token: terraform" \
  -d "{\"tfe-token\":\"$ATLAS_TOKEN\",\"role\":\"$VAULT_ROLE\",\"run-id\":\"$TF_VAR_TFE_RUN_ID\", \"workspace\":\"$TF_VAR_TFC_WORKSPACE_NAME\"}" \
  $VAULT_ADDR/$VAULT_PATH | ./jq -r .auth.client_token)

echo $VAULT_TOKEN > ~/.vault-token
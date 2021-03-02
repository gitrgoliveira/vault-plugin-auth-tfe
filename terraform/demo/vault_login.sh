#!/bin/bash
# Exit if any of the intermediate steps fail
set -e
# jq is not present by default in TFC/TFE
curl -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod a+x jq

eval "$(./jq -r '@sh "export ROLE=\(.role); export VAULT_ADDR=\(.VAULT_ADDR); export VAULT_LOGIN_PATH=\(.VAULT_LOGIN_PATH)"')"

export VAULT_TOKEN=$(curl -X PUT -H "X-Vault-Request: true" \
  -H "X-Vault-Token: terraform" \
  -d "{\"atlas-token\":\"$ATLAS_TOKEN\",\"role\":\"$ROLE\",\"run-id\":\"$TF_VAR_TFE_RUN_ID\", \"workspace\":\"$TF_VAR_TFC_WORKSPACE_NAME\"}" \
  $VAULT_ADDR/$VAULT_LOGIN_PATH | ./jq -r .auth.client_token)

# ./jq --compact-output -n $VAULT_TOKEN
echo "{\"VAULT_TOKEN\": \"${VAULT_TOKEN-unauthorized}\"}"

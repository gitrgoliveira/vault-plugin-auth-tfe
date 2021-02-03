#!/bin/bash
# Exit if any of the intermediate steps fail
set -e
# jq is not present by default in TFC/TFE

curl -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod a+x jq

eval "$(./jq -r '@sh "export ROLE=\(.role); export VAULT_ADDR=\(.VAULT_ADDR)"')"

export VAULT_TOKEN=$(curl -X PUT -H "X-Vault-Request: true" \
  -H "X-Vault-Token: terraform" \
  -d "{\"atlas-token\":\"$ATLAS_TOKEN\",\"role\":\"$ROLE\",\"run-id\":\"$TF_VAR_TFE_RUN_ID\", \"workspace\":\"$TF_VAR_TFC_WORKSPACE_NAME\"}" \
  $VAULT_ADDR/v1/auth/tfe-auth/login | ./jq -r .auth.client_token)

# ./jq --compact-output -n $VAULT_TOKEN

# echo -n $VAULT_TOKEN > /root/.vault-token
echo "{\"VAULT_TOKEN\": \"${VAULT_TOKEN-unauthorized}\"}"

# JSON_OUT=$(curl -s --header "Authorization: Bearer $ATLAS_TOKEN" --header 'Content-Type: application/vnd.api+json' $TF_VAR_ATLAS_ADDRESS/api/v2/account/details)
# curl -s --header "Authorization: Bearer dQ2118vtG0oSlA.atlasv1.zjvSwb8He0WfySs0rmFlzIN3qWdeashKRZA81GPIWJXTySwT2Pu4b6ijMsd2uXo3nxc" --header 'Content-Type: application/vnd.api+json' https://tfe.ric.gcp.hashidemos.io/api/v2/account/details
# dQ2118vtG0oSlA.atlasv1.zjvSwb8He0WfySs0rmFlzIN3qWdeashKRZA81GPIWJXTySwT2Pu4b6ijMsd2uXo3nxc


# echo $JSON_OUT
# echo -n "{\"stuff\": "
# # echo -n "curl -s --header \"Authorization: Bearer $ATLAS_TOKEN\" --header \"Content-Type: application/vnd.api+json\"  --request GET $TF_VAR_ATLAS_ADDRESS/api/v2/account/details"
# echo -n $(curl -s --header "Authorization: Bearer $ATLAS_TOKEN" --header 'Content-Type: application/vnd.api+json' $TF_VAR_ATLAS_ADDRESS/api/v2/account/details)
# echo -n \"}"
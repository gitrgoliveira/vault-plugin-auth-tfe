#!/bin/bash
# Exit if any of the intermediate steps fail
# set -e
# jq is not present by default in TFC/TFE
# jq -n --arg ATLAS_TOKEN "$ATLAS_TOKEN" '{"ATLAS_TOKEN":$ATLAS_TOKEN}'

# echo "{\"ATLAS_TOKEN\": \"${ATLAS_TOKEN:-na}\"}"
curl -s -L -o jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod a+x jq

# curl -X PUT -H "X-Vault-Request: true" \
#   -H "X-Vault-Token: dd" \
#   -d '{"atlas-token":"$ATLAS_TOKEN","role":"workspace_role","run-id":"$TF_VAR_TFE_RUN_ID","workspace":"$TF_VAR_TFC_WORKSPACE_NAME"}' \
#   http://88.97.2.109:8200/v1/auth/tfe-auth/login

curl -X PUT -H "X-Vault-Request: true" \
  -H "X-Vault-Token: dd" \
  -d "{\"atlas-token\":\"$ATLAS_TOKEN\",\"role\":\"$ROLE\",\"run-id\":\"$TF_VAR_TFE_RUN_ID\", \"workspace\":\"$TF_VAR_TFC_WORKSPACE_NAME\"}" \
  http://88.97.2.109:8200/v1/auth/tfe-auth/login



# JSON_OUT=$(curl -s --header "Authorization: Bearer $ATLAS_TOKEN" --header 'Content-Type: application/vnd.api+json' $TF_VAR_ATLAS_ADDRESS/api/v2/account/details)
# curl -s --header "Authorization: Bearer 4JEetYiAhnpbyA.atlasv1.0QK42GxOfVXDvRlAKoidn8uZbjwsbbU3hp8k235LMw4NSSStul7UBIA0BBVFT2vsU80" --header 'Content-Type: application/vnd.api+json' https://tfe.ric.gcp.hashidemos.io/api/v2/account/details
# dQ2118vtG0oSlA.atlasv1.zjvSwb8He0WfySs0rmFlzIN3qWdeashKRZA81GPIWJXTySwT2Pu4b6ijMsd2uXo3nxc


# echo $JSON_OUT
# echo -n "{\"stuff\": "
# # echo -n "curl -s --header \"Authorization: Bearer $ATLAS_TOKEN\" --header \"Content-Type: application/vnd.api+json\"  --request GET $TF_VAR_ATLAS_ADDRESS/api/v2/account/details"
# echo -n $(curl -s --header "Authorization: Bearer $ATLAS_TOKEN" --header 'Content-Type: application/vnd.api+json' $TF_VAR_ATLAS_ADDRESS/api/v2/account/details)
# echo -n \"}"
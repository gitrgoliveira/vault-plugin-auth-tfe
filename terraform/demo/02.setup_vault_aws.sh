#! /bin/bash
#
source helper.sh

vault secrets enable  \
   -default-lease-ttl=120s \
   -max-lease-ttl=240s \
   aws || true

vault write aws/config/root \
   region=us-east-1

vault write aws/roles/deploy \
   role_arns=arn:aws:iam::711129375688:role/ricardo_se_demo \
   credential_type=assumed_role \
   default_sts_ttl=1800 \
   max_sts_ttl=3600

# vault read aws/creds/deploy

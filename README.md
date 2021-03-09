# Vault TFE Auth Plugin

The aim of this Vault authentication plugin is to provide Terraform Cloud or Enterprise with a "window of trust", it can use to retrieve secrets from Vault.

This means you will not need to configure any kind of static secret material for your terraform execution to be able to use Vault.

*TFE Auth* is a test authentication plugin for [HashiCorp Vault](https://www.vaultproject.io/). It is meant for demonstration purposes only and should not assume some kind of official support from HashiCorp.

## TFC/E assumptions

 - TFC/E generates a RUN ID that is unique for that TFE Workspace.

 - Plans or applies are always executed within TFC/E (i.e. remote operations)

 - TFC/E generates a different Atlas token during plan and the apply stages

 - The following environment variables are available:
   - ATLAS_TOKEN
   - TF_VAR_TFE_RUN_ID
   - TF_VAR_TFC_WORKSPACE_NAME

 - The *ATLAS_TOKEN* must have permissions to:
   - [Get the current run ID details](https://www.terraform.io/docs/cloud/api/run.html#get-run-details)
   - [Get the current workspace details](https://www.terraform.io/docs/cloud/api/workspaces.html#show-workspace)
   - [Get the token account details](https://www.terraform.io/docs/cloud/api/account.html#get-your-account-details)

## Vault Authentication conditions
This plugin will issue a token then the following criteria are met:

 - The *ATLAS_TOKEN* provided has the above mentioned permissions
 - The *ATLAS_TOKEN* is a *Service account* token
 - The Run ID provided is in the state "planning" or "applying"
 - The Run ID belongs to the Workspace that is being sent.
 - The Workspace name is in the list of the allowed workspaces for that Role.
 - The Workspace belongs to the TFC/E Organisation configured in the auth nackend

## Usage / Demo

All commands can be run using the provided [Makefile](./Makefile). However, it may be educational to look at the commands to gain a greater understanding of how Vault registers plugins. Using the Makefile will result in running the Vault server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev` server allows you to configure the plugin directory as a flag, and automatically registers plugin binaries in that directory. In production, plugin binaries must be manually registered.

> For the AWS demo, please ensure your AWS credentials have been added to the environment.

This will build the plugin binary and start the Vault dev server:
```bash
# Build TFE Auth plugin and start Vault dev server with plugin automatically registered
$ make
```

Now open a new terminal window and run the following commands:

```bash
# Open a new terminal window and export Vault dev server http address
$ export VAULT_ADDR='http://127.0.0.1:8200'

# Enable the TFE plugin
$ vault auth enable -path=tfe-auth vault-plugin-auth-tfe

# Configure the Authentication backend. By default it points to app.terraform.io
$ vault write auth/tfe-auth/config organization=tfc_org

# Add login roles
$ vault write auth/tfe-auth/role/workspace_role workspaces=* policies=default

```

An example of the above can be seen in [terraform/demo/01.setup_vault.sh](terraform/demo/01.setup_vault.sh)

To login using the tfe auth method, this is the command, but it will not work unless it's run within TFC/E.

```bash
$ vault write auth/tfe-auth/login role=workspace_role \
		workspace=$TFC_WORKSPACE_NAME \
		run-id=$TFC_RUN_ID \
		atlas-token=$ATLAS_TOKEN

```

With terraform, use the code in [terraform/demo/login.tf](terraform/demo/login.tf), which will use the script [terraform/demo/vault_login.sh](terraform/demo/vault_login.sh)

### AWS demo

For the AWS demo, vault needs to be setup with [terraform/demo/02.setup_vault_aws.sh](terraform/demo/02.setup_vault_aws.sh) and the code in here is fairly simple [terraform/demo/aws.tf](terraform/demo/aws.tf)

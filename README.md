# Vault TFE Secrets Plugin

TFE is a test auth plugin for [HashiCorp Vault](https://www.vaultproject.io/). It is meant for demonstration purposes only and should never be used in production.

## Usage

All commands can be run using the provided [Makefile](./Makefile). However, it may be instructive to look at the commands to gain a greater understanding of how Vault registers plugins. Using the Makefile will result in running the Vault server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev` server allows you to configure the plugin directory as a flag, and automatically registers plugin binaries in that directory. In production, plugin binaries must be manually registered.

This will build the plugin binary and start the Vault dev server:

```
# Build TFE Auth plugin and start Vault dev server with plugin automatically registered
$ make
```

Now open a new terminal window and run the following commands:

```
# Open a new terminal window and export Vault dev server http address
$ export VAULT_ADDR='http://127.0.0.1:8200'

# Enable the TFE plugin
$ make enable

# Configure the Authentication backend. By default it points to app.terraform.io
$ vault write auth/tfe-auth/config organization=tfc_org

# Add login roles
$ vault write auth/tfe-auth/role/workspace_role workspaces=* policies=default

```

To login using the tfe auth method:

```
$ vault write auth/tfe-auth/login role=workspace_role \
		workspace=$TFC_WORKSPACE_NAME \
		run-id=$TFC_RUN_ID \
		atlas-token=$ATLAS_TOKEN

```

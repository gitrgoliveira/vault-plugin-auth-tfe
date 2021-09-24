# Vault TFE Auth Plugin

The aim of this Vault authentication plugin is to provide Terraform Cloud or Enterprise with a "window of trust", it can use to retrieve secrets from Vault.

This means you will not need to configure any kind of static secret material for your terraform execution to be able to use Vault.

*TFE Auth* is an authentication plugin for [HashiCorp Vault](https://www.vaultproject.io/), but do not assume some kind of official support from HashiCorp. You should make your own decision whether it is safe or not to use this plugin in your environment.

## Workflow overview
![Workflow overview](images/vault_plugin_workflow.png?raw=true "Workflow overview")

## TFE/TFC assumptions

 - TFE/TFC generates a RUN ID that is unique for that TFE Workspace.

 - Plans or applies are always executed within TFE/TFC (i.e. remote operations)
   - These can also be using terraform remote agents.

 - TFE/TFC generates a different token during plan and the apply stages

 - The following environment variables are available:
   - TF_VAR_TFE_RUN_ID
   - TF_VAR_TFC_WORKSPACE_NAME

 - The TFE/TFC token must have permissions to:
   - [Get the current run ID details](https://www.terraform.io/docs/cloud/api/run.html#get-run-details)
   - [Get the current workspace details](https://www.terraform.io/docs/cloud/api/workspaces.html#show-workspace)
   - [Get the token account details](https://www.terraform.io/docs/cloud/api/account.html#get-your-account-details)

### Retrieving the TFE/TFC token
The TFE/TFC token lives in more than one place. I recommend using the credentials file location.

The credentials file within the TFE/TFC worker lives one of these places, depending if you are using Terraform Agents:
 - `/tmp/cli.tfrc` for code run within TFC/TFE
 - `${path.cwd}/../.terraformrc` for code running in TFC Agents

The TFE/TFC token also exists as an environment variable *ATLAS_TOKEN*. See [terraform/demo/login_env.tf.example](terraform/demo/login_env.tf.example) for an example of that.

## Vault Authentication
### Conditions
This plugin will issue a token then the following criteria are met:

 - The TFE/TFC Token provided has the above mentioned permissions
 - The TFE/TFC Token is a *Service account* token
 - The Run ID provided is in the state "planning" or "applying"
 - The Run ID belongs to the Workspace that is being sent.
 - The Workspace name is in the list of the allowed workspaces for that Role.
 - The Workspace belongs to the TFC/E Organisation configured in the auth backend

### Vault clients / identity
This authentication backend can be configured to use/create different entities depending on the run status: `planning` or `applying`.

This means there is an entity metadata entry key `RunStatus` and the display name is changed to `<organization>/<workspace>/<run_status>`

To achieve this behaviour, you need to set `use_run_status` to `true`
``` bash
vault write auth/tfe-auth/config organization=<org> use_run_status=true
```

With this you can issue **read-only** policies/credentials to entities in `planning` mode and higher priveleged credentials for `applying` mode.
_This is particularly useful when using VCS backed workspaces with speculative plans_.

#### Note
> This setting may increase your pipeline security, but it will also double up on the number of Vault Clients required per TFC/E Workspace


## Usage / Demo

All commands can be run using the provided [Makefile](./Makefile). However, it may be educational to look at the commands to gain a greater understanding of how Vault registers plugins. Using the Makefile will result in running the Vault server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev` server allows you to configure the plugin directory as a flag, and automatically registers plugin binaries in that directory. In production, plugin binaries must be manually registered.

> For the AWS demo, please ensure your AWS credentials have been added to the environment.

This will build the plugin binary and start the Vault dev server:
```bash
# Build TFE Auth plugin and start Vault dev server with plugin automatically registered
$ make
```

A binary can also be downloaded from [the releases page](https://github.com/gitrgoliveira/vault-plugin-auth-tfe/releases).

If instead you are installing this plugin in your Vault environment, then please do this first
```bash
chmod +x /etc/vault.d/plugins/vault-plugin-auth-tfe
# if your vault is using mlock (which it does by default)
setcap cap_ipc_lock=+ep /etc/vault.d/plugins/vault-plugin-auth-tfe
vault plugin register -sha256=84ff2af144de37792d546207e42b618cef027a692bf0d32c324555ba28ca301e auth vault-plugin-auth-tfe
```
The sha256 sum available from [the releases page](https://github.com/gitrgoliveira/vault-plugin-auth-tfe/releases) only validates the zip file.
To get the sha256 sum, please run the following command:
```bash
shasum -a 256 /etc/vault.d/plugins/vault-plugin-auth-tfe > /tmp/vault-plugin-auth-tfe.sha256
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

With terraform, use the code in [terraform/demo/login_file.tf](terraform/demo/login_file.tf)
```
provider "vault" {
  address    = "http://vault_address:8200"
  token_name = "terraform-${var.TFE_RUN_ID}"
  auth_login {
    path = "auth/tfe-auth/login"
    parameters = {
      role      = "workspace_role"
      workspace = var.TFC_WORKSPACE_NAME
      run-id    = var.TFE_RUN_ID
      # For code that is running within TFC/TFE or using an external agent
      tfe-credentials-file = try(filebase64("${path.cwd}/../../../.terraformrc"),
                                  filebase64("/tmp/cli.tfrc"))
    }
  }
}
```

### AWS demo

For the AWS demo, vault needs to be setup with [terraform/demo/02.setup_vault_aws.sh](terraform/demo/02.setup_vault_aws.sh) and the code in here is fairly simple [terraform/demo/aws.tf](terraform/demo/aws.tf)

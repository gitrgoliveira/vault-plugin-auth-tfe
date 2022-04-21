variable "TFC_WORKSPACE_NAME" {
  type        = string
  description = "Workspace of the current run. This is read from the environment and should not be manually set."
}

variable "TFE_RUN_ID" {
  type        = string
  description = "RunID of the current run. This is read from the environment and should not be manually set."
}

variable "VAULT_ADDR" {
  type        = string
  description = "Vault address"
  default = ""
}

variable "VAULT_LOGIN_PATH" {
  type        = string
  description = "Vault login path"
  default = ""
}



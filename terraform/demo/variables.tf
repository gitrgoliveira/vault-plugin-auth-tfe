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
  default     = "http://88.97.2.109:8200"
}


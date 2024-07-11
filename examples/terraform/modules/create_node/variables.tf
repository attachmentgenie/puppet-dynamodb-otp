variable "fqdn" {
  description = "fqdn for node"
  type        = string
}

variable "key_pair_name" {
  type    = string
  default = "otp"
}

variable "puppetserver" {
  type = string
}

variable "security_group_name" {}

variable "subnet_id" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}

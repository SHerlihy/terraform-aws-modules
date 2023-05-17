variable "open_port" {
type = number
default = null
}

variable "subnet_ids" {
  description = "The subnet IDs to deploy to"
  type        = list(string)
}


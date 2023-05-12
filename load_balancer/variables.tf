variable "http_open" {
type = number
default = null
}

variable "subnet_ids" {
  description = "The subnet IDs to deploy to"
  type        = list(string)
}

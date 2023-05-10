variable "http_open" {
type = number
default = 8080
}

variable "lb_security_group_ids" {
type = list(string)
default = ["default"]
}

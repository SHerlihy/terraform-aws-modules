variable "initial_ip" {
  type = string
}

variable "initial_user" {
  default = "ubuntu"
}

variable "initial_pvt_key" {
  type = string
}

variable "accessible_pvt_key_source" {
  type = string
}

variable "accessible_ip" {
  type = string
}

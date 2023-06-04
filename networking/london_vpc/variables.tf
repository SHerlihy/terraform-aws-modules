variable "vpc_cidr_block" {
    type = string
}

variable "pub_cidr_blocks" {
    type = map(string)
    default = {}
}

variable "pvt_cidr_blocks" {
    type = map(string)
    default = {}
}

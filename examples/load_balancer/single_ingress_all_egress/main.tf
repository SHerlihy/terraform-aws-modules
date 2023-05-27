terraform {
      required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
    region = "eu-west-2"
}

module "security_group" {
source = "../../../security_groups/single_ingress_all_egress"

name = "test_security_group"

open_port = var.open_port 
}

module "load_balancer" {
source = "../../../load_balancer"

open_port = var.open_port 

security_group_ids = [module.security_group.module_security_group_id]

subnet_ids = var.subnet_ids
}

output "lb_dns_name" {
    value = module.load_balancer.lb_dns_name
}

output "target_group_arn" {
    value = module.load_balancer.target_group_arn
}

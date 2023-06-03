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

module "london_vpc" {
  source = "../../../networking/london_vpc"

  pub_cidr_blocks = {
    pub1a = "10.0.1.0/24"
    pub1b = "10.0.2.0/24"
  }

  pvt_cidr_blocks = {
    pvt1a = "10.0.3.0/24"
    pvt1b = "10.0.4.0/24"
  }
}

output "pub_subnet_ids_map" {
  description = "pub subnet ids"
  value       = module.london_vpc.pub_subnet_ids_map
}

output "pvt_subnet_ids_map" {
  description = "pvt subnet ids"
  value       = module.london_vpc.pvt_subnet_ids_map
}


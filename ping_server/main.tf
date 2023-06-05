terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

locals {
  cidr_blocks_ping = {
    ping1 = "10.0.1.0/24"
  }
}

module "image_ubuntu" {
  source = "../ubuntu20_id"
}

module "vpc_ping" {
  source = "../networking/london_vpc"

  vpc_cidr_block = "10.0.0.0/16"

  pub_cidr_blocks = local.cidr_blocks_ping
}

module "security_group_ping_ping" {
  source = "../security_groups/ping"

  name = "ping_ping"

  vpc_id = module.vpc_ping.vpc_id

  ingress_cidr_list = var.cidr_blocks_list_ingress
}

data "template_file" "user_data" {
  template = file("${path.module}/../config_server/serve_hello_world.sh")
}

resource "aws_instance" "ping" {
  for_each = module.vpc_ping.pub_subnet_ids_map

  ami           = module.image_ubuntu.ubuntu20_id
  instance_type = "t2.micro"

  subnet_id = each.value
  vpc_security_group_ids = [
    module.security_group_ping_ping.module_security_group_id
  ]

  user_data = data.template_file.user_data.rendered

  tags = {
    Name = each.key
  }
}

output "public_ip" {
  value = aws_instance.ping["ping1"].public_ip
}

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

locals {
  pub_cidr_blocks = {
    pub1a = "10.0.1.0/24"
    pub1b = "10.0.2.0/24"
  }

  pvt_cidr_blocks = {
    pvt1a = "10.0.3.0/24"
    pvt1b = "10.0.4.0/24"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "template_file" "user_data" {
  template = file("../serve_hello_world.sh")
}

module "initial_vpc" {
  source = "../../../networking/london_vpc"

  vpc_cidr_block = "10.0.0.0/16"

  pub_cidr_blocks = local.pub_cidr_blocks

  pvt_cidr_blocks = local.pvt_cidr_blocks
}

module "security_group" {
  source = "../../../security_groups/single_ingress_all_egress"

  name = "test_security_group"

  open_port = 80

  vpc_id = module.initial_vpc.vpc_id
}

resource "aws_instance" "publics" {
  for_each = module.initial_vpc.pub_subnet_ids_map

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  subnet_id              = each.value
  vpc_security_group_ids = [module.security_group.module_security_group_id]

  user_data = data.template_file.user_data.rendered

  tags = {
    Name = each.key
  }
}

resource "aws_instance" "privates" {
  for_each = module.initial_vpc.pvt_subnet_ids_map

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  subnet_id              = each.value
  vpc_security_group_ids = [module.security_group.module_security_group_id]

  user_data = data.template_file.user_data.rendered

  tags = {
    Name = each.key
  }
}

resource "aws_eip" "nat_gateway" {
    vpc = true
}

resource "aws_nat_gateway" "nat_gateway" {
    allocation_id = aws_eip.nat_gateway.id
    subnet_id = module.initial_vpc.pub_subnet_ids_map.pub1a
}

resource "aws_route_table" "pvt" {
    vpc_id = module.initial_vpc.vpc_id

    route {
        cidr_block = "0.0.0.0/0"
        nat_gateway_id = aws_nat_gateway.nat_gateway.id
    }
}

resource "aws_route_table_association" "pvt" {
    subnet_id = module.initial_vpc.pvt_subnet_ids_map.pvt1a
    route_table_id = aws_route_table.pvt.id
}

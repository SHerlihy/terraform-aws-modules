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
  cidr_blocks_pub = {
    pub1a = "10.0.1.0/24"
    pub1b = "10.0.2.0/24"
  }

  cidr_blocks_pvt = {
    pvt1a = "10.0.3.0/24"
    pvt1b = "10.0.4.0/24"
  }

cidr_blocks_ping = {
    test1 = "10.0.1.0/24"
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

module "vpc_ping" {
  source = "../../../networking/london_vpc"

  vpc_cidr_block = "10.0.0.0/16"

  pub_cidr_blocks = local.cidr_blocks_ping
}

module "vpc_app" {
  source = "../../../networking/london_vpc"

  vpc_cidr_block = "10.0.0.0/16"

  pub_cidr_blocks = local.cidr_blocks_pub

  pvt_cidr_blocks = local.cidr_blocks_pvt
}

module "security_group_app_http" {
  source = "../../../security_groups/single_ingress_all_egress"

  name = "app_http"

  open_port = 80

  vpc_id = module.vpc_app.vpc_id
}

module "security_group_app_ssh" {
  source = "../../../security_groups/single_ingress_all_egress"

  name = "app_ssh"

  open_port = 22

  vpc_id = module.vpc_app.vpc_id
}

module "security_group_app_ping" {
  source = "../../../security_groups/ping"

  name = "app_ping"

  vpc_id = module.vpc_app.vpc_id

  ingress_cidr_list = [local.cidr_blocks_pub.pub1a, local.cidr_blocks_pub.pub1b]
}


module "security_group_ping_ssh" {
  source = "../../../security_groups/single_ingress_all_egress"

  name = "ping_ssh"

  open_port = 22

  vpc_id = module.vpc_ping.vpc_id
}

module "security_group_ping_ping" {
  source = "../../../security_groups/ping"

  name = "ping_ping"

  vpc_id = module.vpc_ping.vpc_id

  ingress_cidr_list = [local.cidr_blocks_pub.pub1a, local.cidr_blocks_pvt.pvt1a]
}

resource "aws_key_pair" "ssh-key" {
    key_name = "ssh-key"
    public_key = file("../id_rsa.pub")
}

resource "aws_instance" "pings" {
  for_each = module.vpc_ping.pub_subnet_ids_map

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  subnet_id              = each.value
  vpc_security_group_ids = [
  module.security_group_ping_ssh.module_security_group_id,
  module.security_group_ping_ping.module_security_group_id
  ]

  key_name =  aws_key_pair.ssh-key.key_name 

  tags = {
    Name = each.key
  }
}

//resource "aws_network_interface_sg_attachment" "ping_ssh-test1" {
//  security_group_id    = module.security_group_ping_ssh.module_security_group_id
//  network_interface_id = aws_instance.pings["test1"].primary_network_interface_id
//}

resource "aws_instance" "publics" {
  for_each = module.vpc_app.pub_subnet_ids_map

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  subnet_id              = each.value
  vpc_security_group_ids = [
  module.security_group_app_http.module_security_group_id,
  module.security_group_app_ping.module_security_group_id
  ]

  user_data = data.template_file.user_data.rendered

  key_name =  aws_key_pair.ssh-key.key_name 

  tags = {
    Name = each.key
  }
}

resource "aws_network_interface_sg_attachment" "app_ssh-pub1a" {
  security_group_id    = module.security_group_app_ssh.module_security_group_id
  network_interface_id = aws_instance.publics["pub1a"].primary_network_interface_id
}

resource "aws_instance" "privates" {
  for_each = module.vpc_app.pvt_subnet_ids_map

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  subnet_id              = each.value
  vpc_security_group_ids = [
  module.security_group_app_http.module_security_group_id,
  module.security_group_app_ping.module_security_group_id
  ]

  user_data = data.template_file.user_data.rendered

  key_name =  aws_key_pair.ssh-key.key_name 

  tags = {
    Name = each.key
  }
}

resource "aws_network_interface_sg_attachment" "app_ssh-pvt1a" {
  security_group_id    = module.security_group_app_ssh.module_security_group_id
  network_interface_id = aws_instance.privates["pvt1a"].primary_network_interface_id
}

resource "aws_eip" "nat_gateway" {
  vpc = true
}

resource "aws_nat_gateway" "nat_gateway" {
  allocation_id = aws_eip.nat_gateway.id
  subnet_id     = module.vpc_app.pub_subnet_ids_map.pub1a
}

resource "aws_route_table" "pvt" {
  vpc_id = module.vpc_app.vpc_id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway.id
  }
}

resource "aws_route_table_association" "pvt" {
  subnet_id      = module.vpc_app.pvt_subnet_ids_map.pvt1a
  route_table_id = aws_route_table.pvt.id
}


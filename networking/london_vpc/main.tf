terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_vpc" "main" {
  cidr_block       = var.vpc_cidr_block
  instance_tenancy = "default"

  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "main"
  }
}

resource "aws_subnet" "pub" {
    for_each = var.pub_cidr_blocks

    vpc_id = aws_vpc.main.id
    cidr_block = each.value
    map_public_ip_on_launch = true
}

resource "aws_subnet" "pvt" {
    for_each = var.pvt_cidr_blocks

    vpc_id = aws_vpc.main.id
    cidr_block = each.value
    map_public_ip_on_launch = true
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
}

resource "aws_route_table_association" "publics" {
  for_each       = aws_subnet.pub
  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

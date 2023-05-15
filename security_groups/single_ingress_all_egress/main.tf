terraform {
      required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_security_group" "module_security_group" {
  name = var.name
}

resource "aws_security_group_rule" "ingress_single_port" {
    type = "ingress"
    security_group_id = aws_security_group.module_security_group.id

    from_port   = var.open_port
    to_port     = var.open_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
}


resource "aws_security_group_rule" "egress_all_ports" {
    type = "egress"
    security_group_id = aws_security_group.module_security_group.id

    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
}

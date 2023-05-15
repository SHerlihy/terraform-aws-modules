terraform {
      required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_lb" "lb_1" {
  name               = "http-load-balancer"
  load_balancer_type = "application"
  subnets            = var.subnet_ids 
  // why not var?
  security_groups    = var.security_group_ids
  }

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.lb_1.arn
  port              = var.open_port
  protocol          = "HTTP"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "404: page not found"
      status_code  = 404
    }
  }
}

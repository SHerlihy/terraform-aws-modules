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
  security_groups    = var.security_group_ids
  }

resource "aws_lb_target_group" "front_end" {
  name        = "lb-alb-tg"
  target_type = "instance"
  port        = var.open_port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
}

//resource "aws_lb_target_group_attachment" "tg_attachment_lb" {    
//    target_group_arn = aws_lb_target_group.front_end.arn
//    target_id        = "i-0cbbbbbbbb12f"
//    port             = var.open_port
//}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.lb_1.arn
  port              = var.open_port
  protocol          = "HTTP"

  default_action {
    type = "forward"
    target_group_arn = aws_lb_target_group.front_end.arn
  }
}

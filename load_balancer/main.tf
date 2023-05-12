data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "alb" {
  name = "lb_1_security_group"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]

  }
}


resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.lb_1.arn
  port              = 80
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

// resource "aws_lb_listener_rule" "asg" {
//   listener_arn = aws_lb_listener.http.arn
//   priority     = 100
// 
//   condition {
//     path_pattern {
//       values = ["*"]
//     }
//   }
// 
//   action {
//     type             = "forward"
//     target_group_arn = aws_lb_target_group.asg.arn
//   }
// }

// resource "aws_lb_target_group" "asg" {
//   name     = "asg-health-check"
//   port     = var.http_open
//   protocol = "HTTP"
//   vpc_id   = data.aws_vpc.default.id
// 
//   health_check {
//     path                = "/"
//     protocol            = "HTTP"
//     matcher             = "200"
//     interval            = 15
//     timeout             = 3
//     healthy_threshold   = 2
//     unhealthy_threshold = 2
//   }
// }

resource "aws_lb" "lb_1" {
  name               = "http-load-balancer"
  load_balancer_type = "application"
  subnets            = data.aws_subnets.default.ids
  security_groups    = [aws_security_group.alb.id]
  }

output "load_balancer_dns" {
  value = aws_lb.lb_1.dns_name
}

// output "target_group_arn" {
//     value = aws_lb_target_group.asg.arn
// }

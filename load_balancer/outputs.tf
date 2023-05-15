output "lb_dns_name" {
  value       = aws_lb.lb_1.dns_name
  description = "The domain name of the load balancer"
}

output "alb_http_listener_arn" {
  value       = aws_lb_listener.http.arn
  description = "The ARN of the HTTP listener"
}

output "target_group_arn" {
    value = aws_lb.lb_1.target_group_arn
}

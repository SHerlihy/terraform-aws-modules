output "ping_server-public_DNS" {
  description = "ping server public IP"
  value       = module.ping_server.public_dns
}

output "pub_server-public_DNS" {
  value = aws_instance.publics["pub1a"].public_dns
}

output "pvt_server-net_access-public_DNS" {
  value = aws_instance.privates["pvt1a"].public_dns
}

output "pvt_server-no_access-public_DNS" {
  value = aws_instance.privates["pvt1b"].public_dns
}

output "pvt_server-net_access-private_DNS" {
  value = aws_instance.privates["pvt1a"].private_dns
}

output "pvt_server-no_access-private_DNS" {
  value = aws_instance.privates["pvt1b"].private_dns
}

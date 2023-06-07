output "ping_server-public_ip" {
  description = "ping server public IP"
  value       = module.ping_server.public_ip
}

output "pub_server-public_ip" {
    value = aws_instance.publics["pub1a"].public_ip
}

output "pvt_server-net_access-public_ip" {
    value = aws_instance.privates["pvt1a"].public_ip
}

output "pvt_server-no_access-public_ip" {
    value = aws_instance.privates["pvt1b"].public_ip
}

output "pub_subnet_ids_map" {
    value = {
        for k, sn in aws_subnet.pub : k => sn.id
    }
}

output "pvt_subnet_ids_map" {
    value = {
        for k, sn in aws_subnet.pvt : k => sn.id
    }
}

output "vpc_id" {
  value = aws_vpc.main.id
}

output "route_table_id" {
    value = aws_route_table.public.id
}

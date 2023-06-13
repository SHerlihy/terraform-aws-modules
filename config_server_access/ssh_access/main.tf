terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "terraform_data" "supply_pvt_key" {
  connection {
    type = "ssh"
    port = "22"

    host = var.initial_ip
    user = var.initial_user

    private_key = var.initial_pvt_key

    timeout = "2m"
  }

  provisioner "file" {
    source      = var.accessible_pvt_key_source
destination = "/home/ubuntu/.ssh/id_rsa"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo chmod 400 /home/ubuntu/.ssh/id_rsa",
      "sudo ssh-keyscan -t rsa ${var.accessible_ip} >> /home/ubuntu/.ssh/known_hosts"
    ]
  }
}

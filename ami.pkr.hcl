packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "aws_region" {}
variable "subnet_id" {}
variable "source_ami" {}
variable "demo_account_id" {}
variable "ssh_username" {
  type    = string
  default = "ec2-user"
}


locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

source "amazon-ebs" "aws-linux" {
  access_key      = "${var.aws_access_key}"
  secret_key      = "${var.aws_secret_key}"
  region          = "${var.aws_region}"
  instance_type   = "t2.micro"
  subnet_id       = "${var.subnet_id}"
  source_ami      = "${var.source_ami}"
  ssh_username    = "${var.ssh_username}"
  ami_name        = "csye6225-fall2022-${local.timestamp}"
  ami_description = "Amazon linux 2 AMI for CSYE-6225"
  ami_users       = ["${var.demo_account_id}"]
  launch_block_device_mappings {
    device_name           = "/dev/xvda"
    volume_size           = "8"
    volume_type           = "gp2"
    delete_on_termination = true
  }
}

build {
  name = "web-app-build"
  sources = [
    "source.amazon-ebs.aws-linux"
  ]
  // provisioner "shell-local" {
  //   inline = ["GOOS=linux GOARCH=amd64 go build -o linux-bin/webapp ."]
  // }
  provisioner "shell" {
    inline = [
      "sudo mkdir -p ~/webservice",
      "sudo chown ${var.ssh_username}:${var.ssh_username} ~/webservice",
    ]
  }
  // provisioner "file" {
  //   source      = "./"
  //   destination = "~/webservice"
  // }
  provisioner "file" {
    source      = "linux-bin/webapp"
    destination = "~/webservice/webapp"
  }
  provisioner "file" {
    source      = ".env"
    destination = "~/webservice/.env"
  }
  provisioner "file" {
    source      = "gowebapp.service"
    destination = "~/webservice/gowebapp.service"
  }
  provisioner "shell" {
    inline = [
      "sleep 30",
      "sudo yum update -y",
    ]
  }
  provisioner "shell" {
    inline = [
      "sleep 5",
      "sudo yum install mysql -y",
    ]
  }
  provisioner "shell" {
    inline = [
      "sleep 5",
      "sudo yum install ruby -y",
      "sudo yum install wget -y",
      "cd /home/ec2-user",
      "wget https://aws-codedeploy-us-east-1.s3.us-east-1.amazonaws.com/latest/install",
      "chmod +x ./install",
      "sudo ./install auto",
      "sudo service codedeploy-agent status",
    ]
  }
  provisioner "shell" {
    inline = [
      "cd ~/webservice",
      "sudo mv gowebapp.service /lib/systemd/system/gowebapp.service",
      "sudo systemctl start gowebapp",
      "sudo systemctl enable gowebapp",
      "sudo systemctl status gowebapp",
    ]
  }
  // provisioner "shell" {
  //   inline = [
  //     "sleep 5",
  //     "sudo yum install mariadb-server -y",
  //     "sudo systemctl start mariadb",
  //     "sudo systemctl enable mariadb",
  //     "echo 'create SCHEMA webservicedb;' | sudo mysql",
  //     "mysqladmin -u root password p@ssword",
  //   ]
  // }
  // provisioner "shell" {
  //   inline = [
  //     "sleep 5",
  //     "sudo yum install golang -y",
  //   ]
  // }
  // provisioner "shell" {
  //   inline = [
  //     "cd ~/webservice",
  //     "go build -o webapp .",
  //     "sudo mv gowebapp.service /lib/systemd/system/gowebapp.service",
  //     "sudo systemctl start gowebapp",
  //     "sudo systemctl enable gowebapp",
  //   ]
  // }
}

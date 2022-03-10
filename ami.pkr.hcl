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


build {
  name = "web-app-build"
  sources = [
    "source.amazon-ebs.aws-linux"
  ]
  provisioner "shell" {
    inline = [
      "echo creating required directories",
      "sleep 5",
      "sudo mkdir -p ~/webservice",
      "sudo chown ${var.ssh_username}:${var.ssh_username} ~/webservice",
    ]
  }
  provisioner "file" {
    source      = "./"
    destination = "~/webservice"
  }
  provisioner "shell" {
    environment_vars = [
      "FOO=foo",
    ]
    inline = [
      "echo updating packages",
      "sleep 20",
      "sudo yum update -y",
    ]
  }
  provisioner "shell" {
    inline = [
      "echo installing mysql",
      "sleep 5",
      "sudo yum install mariadb-server -y",
      "sudo systemctl start mariadb",
      "sudo systemctl enable mariadb",
      "echo 'create SCHEMA webservicedb;' | sudo mysql",
      "mysqladmin -u root password p@ssword",
    ]
  }
  provisioner "shell" {
    inline = [
      "echo installing golang",
      "sleep 5",
      "sudo yum install golang -y",
    ]
  }
  provisioner "shell" {
    inline = [
      "echo running the web app",
      "sleep 5",
      "cd ~/webservice",
      "go build -o webapp .",
      "sudo mv gowebapp.service /lib/systemd/system/gowebapp.service",
      "sudo systemctl start gowebapp.service",
      "sudo systemctl enable gowebapp.service",
    ]
  }
}

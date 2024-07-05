data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "aws_dynamodb_table" "puppet-dynamodb-otp" {
  name = "puppet-dynamodb-otp"
}

data "aws_key_pair" "otp" {
  key_name = "otp"
}

data "aws_security_group" "otp-ssh" {
  name = var.security_group
}

data "aws_subnet" "otp-public" {
  cidr_block = "10.0.1.0/24"
}

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

module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.6.1"

  name = "puppetserver"

  instance_type          = "t3.micro"
  ami                    = data.aws_ami.ubuntu.image_id
  key_name               = module.otp_key.key_pair_name
  vpc_security_group_ids = [module.otp_sg.security_group_id]
  subnet_id              = element(module.otp_vpc.public_subnets, 0)

  iam_instance_profile = module.otp_node_role.iam_instance_profile_name

  tags = local.tags
}

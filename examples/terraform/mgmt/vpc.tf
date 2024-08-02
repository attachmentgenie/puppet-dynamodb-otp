module "otp_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.8.1"

  name = "otp"
  cidr = "10.0.0.0/16"

  azs            = data.aws_availability_zones.available.zone_ids
  public_subnets = ["10.0.1.0/24"]

  tags = local.tags
}

module "otp_sg" {
  source  = "terraform-aws-modules/security-group/aws//modules/ssh"
  version = "5.1.2"

  name   = "otp-ssh"
  vpc_id = module.otp_vpc.vpc_id

  ingress_cidr_blocks = ["0.0.0.0/0"]

  tags = local.tags
}

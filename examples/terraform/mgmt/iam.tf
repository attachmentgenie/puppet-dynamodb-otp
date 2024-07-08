module "otp_key" {
  source  = "terraform-aws-modules/key-pair/aws"
  version = "2.0.3"

  key_name           = "otp"
  create_private_key = true
}

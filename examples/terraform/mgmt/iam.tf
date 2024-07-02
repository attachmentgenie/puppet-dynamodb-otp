module "otp_key" {
  source  = "terraform-aws-modules/key-pair/aws"
  version = "2.0.3"

  key_name           = "otp"
  create_private_key = true
}

module "otp_node_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"
  version = "5.39.1"

  trusted_role_services = [
    "ec2.amazonaws.com",
  ]
  custom_role_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"
  ]
  number_of_custom_role_policy_arns = 1

  allow_self_assume_role  = true
  create_role             = true
  create_instance_profile = true

  role_name = "otp"

  tags = local.tags
}

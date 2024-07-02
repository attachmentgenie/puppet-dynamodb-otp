module "puppet-dynamodb-otp" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "4.0.1"

  billing_mode       = "PAY_PER_REQUEST"
  name               = "puppet-dynamodb-otp"
  hash_key           = "fqdn"
  ttl_enabled        = true
  ttl_attribute_name = "expire_at_unix"

  attributes = [
    {
      name = "fqdn"
      type = "S"
    }
  ]

  tags = local.tags
}
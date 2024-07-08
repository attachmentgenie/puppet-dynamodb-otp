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
    },
    {
      name = "expire_at_unix"
      type = "N"
    },
    {
      name = "otp_token"
      type = "S"
    }
  ]

  global_secondary_indexes = [
    {
      name               = "expire_at_unix"
      hash_key           = "expire_at_unix"
      range_key          = "otp_token"
      projection_type    = "INCLUDE"
      non_key_attributes = ["otp_token"]
    }
  ]

  tags = local.tags
}
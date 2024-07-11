module "clients" {
  source = "../modules/create_node"

  for_each = toset(local.clients)
  fqdn     = "${each.key}.${local.domain}"

  puppetserver        = var.puppetserver
  security_group_name = var.security_group
  subnet_id           = data.aws_subnet.otp-public.id

  tags = local.tags
}
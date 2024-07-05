output "private_key" {
  value     = trimspace(module.otp_key.private_key_openssh)
  sensitive = true
}

output "security_group" {
  value = module.otp_sg.security_group_name
}

output "puppetserver" {
  value = module.puppetserver.private_ip
}

output "public_ip" {
  value = module.puppetserver.public_dns
}
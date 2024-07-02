output "foo" {
  value = "bar"
}

# output "private_key" {
#   value     = tls_private_key.ssh_key.private_key_pem
#   sensitive = true
# }

# output "public_key" {
#   value     = tls_private_key.ssh_key.public_key_openssh
#   sensitive = true
# }

# output "puppetserver" {
#   value = aws_instance.puppetserver.public_dns
# }
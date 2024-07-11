output "public_dns" {
  value = tomap({ for i, node in module.clients : i => node.public_dns })
}
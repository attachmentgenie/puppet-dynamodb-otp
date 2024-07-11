resource "random_string" "token" {
  length  = 16
  special = false
}

module "otp_token" {
  source = "../puppet-dynamodb-otp//modules/token_table_item"

  fqdn  = var.fqdn
  token = random_string.token.result
}

module "node" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.6.1"

  name = var.fqdn

  ami                    = data.aws_ami.ubuntu.image_id
  key_name               = data.aws_key_pair.otp.key_name
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [data.aws_security_group.otp-ssh.id]

  associate_public_ip_address = true
  instance_type               = "t3.micro"
  user_data                   = <<-EOT
#cloud-config
manage_etc_hosts: false
fqdn: ${var.fqdn}
prefer_fqdn_over_hostname: true
write_files:
  - path: /etc/puppetlabs/puppet/csr_attributes.yaml
    owner: root:root
    permissions: '0755'
    content: |
      custom_attributes:
        challengePassword: ${random_string.token.result}
  - path: /var/cache/configure-puppet.sh
    owner: root:root
    permissions: '0755'
    content: |
      #!/bin/sh
      #
      # Script body start
      wget -qO - https://raw.githubusercontent.com/puppetlabs/install-puppet/main/install.sh | bash -s -- -c puppet8
      /opt/puppetlabs/bin/puppet config set --section agent environment production
      /opt/puppetlabs/bin/puppet config set --section main splay true
      /opt/puppetlabs/bin/puppet config set --section main splaylimit 300
      /opt/puppetlabs/bin/puppet config set --section main runinterval 300
      systemctl enable puppet
      systemctl start puppet
      # Script body end
runcmd:
  - systemctl disable ufw
  - echo '${var.puppetserver} puppet' >> /etc/hosts
  - /var/cache/configure-puppet.sh
  EOT

  tags = var.tags
}
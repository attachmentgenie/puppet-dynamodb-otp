resource "random_string" "token" {
  length  = 16
  special = false
}

resource "time_static" "current_date" {}

resource "aws_dynamodb_table_item" "token" {
  table_name = data.aws_dynamodb_table.puppet-dynamodb-otp.name
  hash_key   = data.aws_dynamodb_table.puppet-dynamodb-otp.hash_key

  item = <<ITEM
{
  "expire_at_unix": {"N": "${time_static.current_date.unix + 300}"},
  "fqdn": {"S": "client1.ec2.internal"},
  "otp_token": {"S": "${random_string.token.result}"}
}
ITEM
}

module "client1" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.6.1"

  name = "client1"

  ami                    = data.aws_ami.ubuntu.image_id
  key_name               = data.aws_key_pair.otp.key_name
  subnet_id              = data.aws_subnet.otp-public.id
  vpc_security_group_ids = [data.aws_security_group.otp-ssh.id]

  associate_public_ip_address = true
  instance_type               = "t3.micro"
  user_data                   = <<-EOT
#cloud-config
manage_etc_hosts: false
fqdn: client1
prefer_fqdn_over_hostname: true
hostname: client1
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

  tags = local.tags
}
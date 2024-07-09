data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

module "puppetserver" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "5.6.1"

  name = "puppetserver"

  ami                    = data.aws_ami.ubuntu.image_id
  key_name               = module.otp_key.key_pair_name
  subnet_id              = element(module.otp_vpc.public_subnets, 0)
  vpc_security_group_ids = [module.otp_sg.security_group_id]

  associate_public_ip_address = true

  create_iam_instance_profile = true
  iam_role_description        = "IAM role for puppetserver"
  iam_role_policies = {
    AmazonDynamoDBFullAccess = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
  }

  instance_type = "t3.medium"
  user_data     = <<-EOT
## template: jinja
#cloud-config
fqdn: puppet
prefer_fqdn_over_hostname: true
hostname: puppet
write_files:
  - path: "/etc/environment"
    append: true
    content: |
      export AWS_REGION="{{ v1.region }}"
  - path: /var/cache/configure-puppetserver.sh
    owner: root:root
    permissions: '0755'
    content: |
      #!/bin/sh
      #
      # Script body start
      wget -qO - https://raw.githubusercontent.com/puppetlabs/install-puppet/main/install.sh | bash -s -- -c puppet8
      apt install -y puppetserver
      /opt/puppetlabs/bin/puppet config set --section agent environment production
      /opt/puppetlabs/bin/puppet config set --section server autosign /usr/local/bin/puppet-dynamodb-otp
      # Script body end
  - path: "/etc/systemd/system/puppetserver.service.d/env.conf"
    content: |
      [Service]
      Environment="AWS_REGION={{ v1.region }}"
runcmd:
  - systemctl disable ufw
  - /var/cache/configure-puppetserver.sh
  - systemctl enable puppetserver
  - systemctl start puppetserver
  EOT

  tags = local.tags
}

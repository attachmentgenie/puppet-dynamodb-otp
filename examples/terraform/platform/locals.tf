locals {
  clients = ["client1"]
  domain  = "ec2.internal"
  tags = {
    terraform = "true"
    project   = "otp-platform"
  }
}
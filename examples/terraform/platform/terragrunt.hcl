include "root" {
  path = find_in_parent_folders()
}

dependency "mgmt" {
  config_path = "../mgmt"
}

inputs = {
  puppetserver   = dependency.mgmt.outputs.puppetserver
  security_group = dependency.mgmt.outputs.security_group
}

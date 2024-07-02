include "root" {
  path = find_in_parent_folders()
}

dependency "mgmt" {
  config_path = "../mgmt"
}

inputs = {
  foo = dependency.mgmt.outputs.foo
}

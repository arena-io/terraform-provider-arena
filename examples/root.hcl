generate "provider" {
  path = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents = <<EOF
terraform {
  required_providers {
    arena = {
      source  = "arena-io/arena"
    }
  }
}

provider "arena" {
  server_url = "http://localhost:18080/api/v1"
}

EOF
}
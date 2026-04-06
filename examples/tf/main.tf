# Copyright (c) ArenaML Labs Pvt Ltd.

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


resource "arena_cluster_manager" "def" {
  name = "tf_test_engine"
  kind = "nomad"
  spec = file("${path.module}/default-engine-spec.json")
}

resource "arena_org" "nav_dev" {
  name = "nav-dev"
  description = "developer group for navigation systems"
}

resource "arena_team" "auto_flight" {
  name   = "path-planning"
  role = "devs"
  org_id = arena_org.nav_dev.id
  description = "dev team working on autonomous path planning"
  config = {
    allow_cross_orgs = true
  }
}

resource "arena_user" "rock_star" {
  email = "rock@star.univ"
  name  = "rock star"
}

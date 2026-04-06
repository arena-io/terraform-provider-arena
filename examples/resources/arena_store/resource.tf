
resource "arena_store" "ais_staging" {
  name = "ais-staging"
  kind = "aistore"
  basepath = "/arena-ml"
  endpoint = "http://10.1.1.2:51080"
  org_id = "5ca1ab1e-0000-4000-a000-000000000000"
  config = {
    auth = jsonencode({
      token = "top-secret"
    })
    max_objects = 1e6
    capacity_gb = 1000
  }
}
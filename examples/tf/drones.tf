
resource "arena_drone_profile" "rx" {
  description = "rpi-2"
  kind        = "sbc"
  name        = "dfrobota"
  spec = {
    arch = "arm64"
    memory_in_gb = 32

    compute = {
      "gpu"  = 2
      "cuda" = 13
    }
    storage = {
      main = {
        guid     = "1234"
        dev_path = "/mnt/dev"
        capacity = 9
      }
    }
  }
}

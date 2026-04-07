resource "arena_pipeline" "test_pln" {
  name        = "some-pl"
  description = "some pipeline"
}

resource "arena_pipeline_input" "input_a" {
  name        = "in-A"
  pipeline_id = arena_pipeline.test_pln.id
}


resource "arena_pipeline_step" "step_a" {
  pipeline_id = arena_pipeline.test_pln.id
  name        = "step-a"
  kind        = "docker"
  config = {
    image = "debian:latest"
    run_spec = jsonencode({
      cpu    = 2000
      memory = 4096
    })
  }
}

resource "arena_pipeline_output" "out_b" {
  name        = "out-a"
  pipeline_id = arena_pipeline.test_pln.id

}


resource "arena_pipeline_dag" "test_pl_dag" {
  pipeline_id = arena_pipeline.test_pln.id

  input_edges = [
    {
      node_id    = arena_pipeline_input.input_a.id
      from_bases = [arena_basis.test_basis.id]
      to_steps   = [arena_pipeline_step.step_a.id]
      # from_inputs = []
      # to_inputs = []
    }
  ]

  output_edges = [
    {
      node_id   = arena_pipeline_output.out_b.id
      from_step = arena_pipeline_step.step_a.id
      # to_inputs = []
    }
  ]
}

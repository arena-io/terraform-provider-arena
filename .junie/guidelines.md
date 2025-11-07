Project-specific development guidelines (Last verified: 2025-10-21 14:38 local)

Audience: Advanced Go/Terraform provider developers working on this repo.

1. Build and configuration
- Go/toolchain
  - go.mod pins: go 1.25.1 with toolchain go1.25.1. Use Go ≥ 1.25 toolchain for consistency.
  - mise and Task are configured for tooling: see mise.toml and taskfile.yml.
- Compile provider binary
  - Quick build: task build-tf (wraps: go build -o ./bin/terraform-provider-arenaml)
  - Manual: go build -o ./bin/terraform-provider-arenaml .
  - Minimal example exists at `tf/main.tf` .
- Generated code and codegen
  - The generator/schema/ folder is produced via HashiCorp’s tfplugingen-framework from specs/*.json. Task targets:
    - task gen-tf-schema  # refreshes data-sources and resources into generator/schema. 
    - task gen-rest-client (in generator/client)  # go generate for the REST client
  - Note: As of this snapshot, generator/schema contains duplicate symbol definitions for basis schemas that will not compile as part of the workspace. These files are intended as scaffolding/reference for manual curation; do not import these packages into the provider until they’re reconciled. See “Testing” section for how to exclude them.
  - the code in generator/schema needs correction before they can be used. 
  - working and actually used schema is in `internal/schama`
  
- Docs generation
  - If you update schemas/provider, consider regenerating docs with tfplugindocs once the compiled provider exposes the schema. The repo already includes docs/*.
- Linting and hooks
  - Run: task lint (golangci-lint). Pre-commit hooks are configured; initialize with task setup-pre-commit.

2. Testing
- Current status
  - Helper package tests compile and pass. The generated schema packages under generator/schema/*basis* fail to build due to duplicate type/func declarations (intentional/unresolved scaffolding), so whole-repo test invocation must exclude those directories.
- Run tests (package-scoped)
  - Fast path: go test ./helper
- Run tests (all non-generated packages)
  - POSIX shells: go test $(go list ./... | grep -v '^.*/generator/schema/')
  - Using Task with a filter: task test -- <Regex> (wraps go test -run '<Regex>' ./...). Be aware this still traverses ./...; prefer the go list + grep form to exclude generator/ until it compiles.
- Adding a new test
  - Place *_test.go files alongside the code under test, using the standard Go testing package.
  - Prefer table-driven tests and avoid cross-package imports into generator/schema until that code compiles cleanly.
- Verified example (created and executed during this session)
  - A minimal smoke test was added temporarily under helper and executed with:
    go test ./helper
    Result: ok (pass). The file was then removed to keep the tree clean, per the request.

3. Additional development information
- Provider layout
  - Provider entry: internal/provider/provider.go; resources/data sources currently implemented for “engine” under internal/provider. top-level main.go wires up the framework provider.
  - Public helper routines: helper/json.go, helper/tags.go. These are safe places for shared transformation logic.
- Known pitfalls
  - generator/schema duplication: Running gen-tf-schema can produce duplicate symbols (e.g., RunSpecType, NewRunSpecValue*, etc.) for basis resources/data sources. Treat output as a starting point and prune/merge types before including them in build targets.
  - Testing ./... will fail until generator/schema is reconciled; see commands above to exclude it.
  - The helper/json.go type-conversion utility is strict about type matches and logs type mismatches via tflog. When extending schemas, ensure json tags align with tfsdk tags to avoid silent no-ops in reflection-based mapping.
- Style and tooling
  - Use golangci-lint (task lint). Keep imports tidy and leverage pre-commit.
  - Logging: use github.com/hashicorp/terraform-plugin-log/tflog; avoid fmt.Print in provider code paths.
  - Framework: github.com/hashicorp/terraform-plugin-framework v1.15.0 is in use. Prefer basetypes and jsontypes helpers for complex JSON attributes to keep state/model parity.
- Working with examples
  - An example provider config is in `tf/main.tf` (source = "arenaml/arenaml", version pinned). For local dev overrides, omit the version or allow it but ensure the dev override takes precedence.
  - A minimal Terraform configuration for the engine resource/data source is in tf/; it’s useful for manual ad-hoc E2E once the provider binary is built and credentials are configured. Ensure you point Terraform to the locally built plugin as noted above.
- Environment and secrets
  - The Taskfile loads some optional dotenv files (.env, ~/.secrets/ctrl-dev.env, ~/.arenaml/ctrl-dev.env). If you rely on those for credentials or API endpoints, document the variables in README/USAGE and avoid checking secrets into VCS.

4. Quick-reference commands
- Build provider: task build-tf
- Lint: task lint
- Generate TF schema scaffolding: task gen-tf-schema
- Run unit tests for helper: go test ./helper
- Run unit tests excluding generator/schema: go test $(go list ./... | grep -v '^.*/generator/schema/')

Notes
- All test commands provided above were verified on 2025-10-21. The full-repo go test ./... currently fails due to generated schema duplication; this is expected until generator outputs are reconciled.
- This document intentionally avoids basic Go/Terraform explanations and focuses on idiosyncrasies of this repo.
- don't do any validation like building the binary or running tests
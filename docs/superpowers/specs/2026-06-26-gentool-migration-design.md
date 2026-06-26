# Gentool Migration Design

## Context

The project currently generates GORM DAO code through `cmd/gen/generate.go`, a custom Go entrypoint that imports `gorm.io/gen`, loads `UserRPC.Mysql` from `config.yaml`, and writes generated code into `apps/user-rpc/internal/dao/query` plus the adjacent `model` package.

The existing generator targets four tables:

- `user`
- `student`
- `college`
- `mini_program_user`

The generated code is consumed through stable package paths:

- `github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model`
- `github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query`

Current consumers rely on `query.Use(db)`, `query.Query`, and generated table-specific query fields such as `Query.User` and `Query.Student`.

## Goal

Migrate model/query generation from the hand-written `cmd/gen` entrypoint to the official `gorm.io/gen/tools/gentool` CLI while preserving the existing generated-code package layout and runtime API used by the service and repository layers.

## Non-Goals

- Do not change API or RPC code generation.
- Do not refactor repository or logic packages unless generated API compatibility requires it.
- Do not modify unrelated existing user-api changes in the working tree.
- Do not broaden table generation beyond the four current tables.

## Constraints

`gentool` is a separate module from the project's current `gorm.io/gen` dependency. The project currently depends on `gorm.io/gen v0.3.27`; the latest listed library version is `v0.3.28`, while the locally cached CLI module is `gorm.io/gen/tools/gentool@v0.0.2`.

`gentool` supports DSN, database type, table list, output path, output file, model package path, and standard field switches. It does not expose the full `GenerateModel` option surface used by the current custom generator, specifically:

- custom data type mapping such as `tinyint -> int8`
- custom field type override for `deleted_at`
- custom GORM tag mutation for soft delete
- custom JSON tag override for `deleted_at`

The current schema in `deploy/mysql/ddl.sql` does not define a `deleted_at` column, so the soft-delete options are currently defensive configuration rather than active schema behavior.

## Recommended Approach

Use `gentool` as the generation interface and keep a small project-owned wrapper only where the CLI cannot know project-specific configuration.

Implementation should:

1. Add a reproducible way to install or invoke `gorm.io/gen/tools/gentool` at a fixed version.
2. Add a checked-in `gen.yml` or Makefile target that records the four-table table list and output path.
3. Keep generated output paths unchanged:
   - query output: `apps/user-rpc/internal/dao/query`
   - model output: `apps/user-rpc/internal/dao/model`
4. Remove `cmd/gen/generate.go` if the CLI path can generate compatible code without custom Go logic.
5. If `gentool` cannot preserve required generated API compatibility, keep only the smallest compatibility wrapper needed and document that `gentool` remains the externally invoked tool.
6. Preserve `make generate-model` as the developer-facing command.

## Alternatives Considered

### Pure CLI Invocation

Run `gentool` directly from `make generate-model` with DSN, tables, and output path flags.

This is the cleanest migration and should be attempted first. Its risk is that CLI flags cannot represent all existing generator customizations.

### Gentool Config File

Check in a `gen.yml` containing database type, table list, output path, and field switches, then call `gentool -c <path>`.

This is more readable than a long Makefile command and makes generation settings easier to review. It still cannot express custom type maps or field tag mutations.

### Thin Compatibility Wrapper

Keep a small Go wrapper that shells out to or configures gentool-compatible generation while preserving project-specific behavior.

This is the fallback if direct `gentool` output breaks existing imports, field names, or query API shape. It is less pure, but it keeps migration safe.

## Data Flow

Developer runs `make generate-model`.

The Makefile invokes `gentool` at the pinned version. The command reads the MySQL DSN from the chosen config source or receives it as a variable, connects to MySQL, introspects the four configured tables, and regenerates model/query files under `apps/user-rpc/internal/dao`.

Application code continues to open the database through existing service context code and call `query.Use(db)` as before.

## Error Handling

Generation should fail fast when:

- no DSN is supplied
- the database cannot be reached
- any configured table cannot be introspected
- generated code does not compile

The Makefile target should produce a non-zero exit code in each case so CI and local development both catch failures.

## Testing and Verification

Implementation must verify:

- `make generate-model` invokes gentool successfully in an environment with a valid DSN.
- `go test ./...` still compiles existing generated-code consumers.
- The resulting git diff is limited to intended generation tooling and generated DAO output.

If a live database is unavailable during verification, run the strongest available checks and report the gap explicitly.

## Open Decisions

During implementation, choose between direct CLI flags and `gen.yml` after checking which produces the clearest and most reproducible command in this repository.

The expected default is a checked-in config file for static generation settings plus a Makefile variable for the DSN, so secrets do not enter the repository.

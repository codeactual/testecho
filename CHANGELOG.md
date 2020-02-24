# Change Log

## v0.1.4

> This release updates several first/third-party dependencies.

- feat
  - --version now prints details about the build's paths and modules.
- notable dependency changes
  - Bump github.com/pkg/errors to v0.9.1.
  - Bump internal/cage/... to latest from monorepo.
- refactor
  - Migrate to latest cage/cli/handler API (e.g. handler.Session and handler.Input) and conventions (e.g. "func NewCommand").

## v0.1.3

- dep: update first-party dependencies under `internal`

## v0.1.2

- feat: `NewCmdArgs` and `NewCmdString` no longer require an argument
- refactor: remove unused fixtures from `internal/cage/`

## v0.1.1

- feat: add `testecho` package
- dep: update `vendor` and first-party dependencies under `internal`

## v0.1.0

- feat: initial project export from private monorepo

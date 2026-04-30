# Changelog

All notable changes to this project are documented here.

## [Unreleased]

## [2.0.0] - 2026-04-30

### Changed
- Switched CLI to command mode: `add`, `remove`, `list`, `sync`, `help`, `version`.
- Registry is the source of truth for the current icon set.
- Generation now outputs `icons.templ` and `registry.templ`.

### Removed
- Legacy flag-based CLI flow.
- Legacy grouping-based generation and extra output file generation.
- Dry-run and verbose command options from the CLI.

# Lucide Templ Generator

A focused CLI for managing a persistent Lucide icon set and generating Templ components.

## Install

```bash
go install github.com/peterszarvas94/lucide-templ-gen/cmd/lucide-gen@latest
```

## Usage

```text
Usage:
  lucide-gen <command> [arg]
Commands:
  add <icons>       Add icon(s) to current set
  add --all         Add all Lucide icons
  remove <icons>    Remove icon(s) from current set
  remove --all      Remove all icons
  list              List current icons
  sync              Regenerate files from current icon set
  help              Show help
  version           Show version
Arg:
  <icons>           Comma-separated icon names
                    Example: "arrow-left,arrow-right"
Options:
  --output <dir>    Output folder (default: ./icons)
```

> [!NOTE]
> This version is not compatible with [v1](https://github.com/peterszarvas94/lucide-templ-gen/releases/tag/v1.3.3) or the [upstream repository](https://github.com/riclib/lucide-templ-gen).

## Quick Start

```bash
# Start with an empty set
lucide-gen list

# Add icons
lucide-gen add "arrow-left,arrow-right"

# See current set
lucide-gen list

# Remove one icon
lucide-gen remove "arrow-left"

# Regenerate from current registry set
lucide-gen sync
```

## The registry

- The current icon set is read from `registry.templ` in `--output`.
- `add` and `remove` update that set, then regenerate files.
- `sync` regenerates from the registry set only.
- If `registry.templ` is malformed, commands fail and stop.

## Generated Files

By default (`--output ./icons`), generation writes:

- `icons.templ`
- `registry.templ`

## Testing

Run fast unit tests:

```bash
go test ./...
```

Run CLI integration test (real temp folder + real command execution):

```bash
go test -tags=integration ./cmd/lucide-gen -run TestCLIWorkflowIntegration -count=1
```

## License

MIT. See `LICENSE`.

Upstream license: `licenses/riclib/lucide-templ-gen/LICENSE`.

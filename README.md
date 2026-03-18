# otelext
[![CI](https://github.com/thumbrise/otelext/actions/workflows/ci.yml/badge.svg)](https://github.com/thumbrise/otelext/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/thumbrise/otelext.svg)](https://pkg.go.dev/github.com/thumbrise/otelext)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
Extensions for [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/) — reusable filters and samplers for fine-grained telemetry control.
## Features
- **Filter registry** — global, thread-safe registry of attribute-based filters (`signal.Filter` interface)
- **FilterBasedSampler** — trace sampler that drops spans based on registered filters
- **CompositeSampler** — combines multiple `sdktrace.Sampler` instances with a restrict-only strategy (any sampler can restrict the decision, but never upgrade it)
## Installation
```bash
go get github.com/thumbrise/otelext
```
Requires **Go 1.25+** and OpenTelemetry SDK **v1.42+**.
## Quick Start
### Register a filter and use FilterBasedSampler
```go
package main
import (
	"github.com/thumbrise/otelext/signal"
	"github.com/thumbrise/otelext/signal/trace"
)
func main() {
	// Register your custom filter (must implement signal.Filter)
	signal.RegisterFilter(myFilter)
	// Create a sampler that drops spans matched by registered filters
	sampler := trace.NewFilterBasedSampler()
	// Use sampler in your TracerProvider setup
	// sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler))
}
```
### Compose multiple samplers
```go
sampler := trace.NewCompositeSampler(
	trace.NewFilterBasedSampler(),
	sdktrace.TraceIDRatioBased(0.5),
)
// The strictest decision wins: if any sampler returns Drop, the span is dropped.
```
## API
### `signal.Filter` interface
```go
type Filter interface {
	ShouldDrop(ctx context.Context, attrs attribute.Set) bool
	Key(ctx context.Context) string
	Description(ctx context.Context) string
}
```
### `signal` package
| Function                | Description                                   |
|-------------------------|-----------------------------------------------|
| `RegisterFilter(f)`    | Add a filter to the global registry           |
| `RegisteredFilters()`  | Get a copy of all registered filters          |
| `ClearFilters()`       | Remove all registered filters                 |
### `signal/trace` package
| Type                  | Description                                                        |
|-----------------------|--------------------------------------------------------------------|
| `FilterBasedSampler`  | Sampler that drops spans when any registered filter matches        |
| `CompositeSampler`    | Combines samplers; strictest decision wins (restrict-only strategy)|

## TODO
### Documentation
- [ ] Add package-level GoDoc comments (`signal/doc.go`, `signal/trace/doc.go`)
- [ ] Add `CONTRIBUTING.md` (contribution guide, commit style, PR process)
- [ ] Add `CHANGELOG.md` or set up auto-generation via semantic-release
### Functionality
- [ ] Add built-in `signal.Filter` implementations (by span name, by attributes, regex-based, etc.)
- [ ] Support filters for metrics and logs (currently traces only)
- [ ] Add `Option` pattern for `FilterBasedSampler` (fallback decision, custom logger)
- [ ] Add `Option` pattern for `CompositeSampler` (strategy: restrict-only / permissive / majority)
- [ ] Consider `UnregisterFilter(key)` for dynamic filter management
- [ ] Rework registry. Replace global package level registry with a struct for multi callers possibility
### Tests
- [ ] Add benchmarks for `CompositeSampler` with a large number of samplers
- [ ] Add concurrency tests (parallel filter registration/read)
- [ ] Add integration tests with a real `TracerProvider`
### CI/CD
- [ ] Add CI step for test coverage reporting (codecov / coveralls)
- [ ] Add coverage badge to README
- [x] Set up automatic tag releases
### Infrastructure
- [ ] Add `CODEOWNERS` file
- [ ] Add `SECURITY.md` (responsible disclosure policy)
- [ ] Publish on pkg.go.dev with examples (Example tests)
- [ ] Add `.goreleaser.yml` if binary releases are planned

## Development
### Prerequisites
- [Go 1.25+](https://go.dev/dl/)
- [Task](https://taskfile.dev/) (task runner)
- [golangci-lint v2.4+](https://golangci-lint.run/)
### Commands
```bash
task lint      # Run golangci-lint + license header checks
task test      # Run tests + benchmarks
task license   # Auto-fix missing license headers
```
## License
[Apache License 2.0](LICENSE) — Copyright 2026 thumbrise
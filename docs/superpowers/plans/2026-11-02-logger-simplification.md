# Logger Simplification Implementation Plan

> **For agentic workers:** Execute the approved logging refactor task by task, with regression tests for observable logging behavior.

**Goal:** Use native zap logging while keeping project configuration, context enrichment, and OpenTelemetry/GORM integration.

**Architecture:** A single `New(*Options) (*zap.Logger, error)` builds the logger without global side effects. `Init(*Options) error` installs it with `zap.ReplaceGlobals` and redirects standard logging. A shared level filter and sampler wrap the combined local/OTel core. `L(context.Context) *zap.Logger` adds request fields and both local and OTLP trace correlation.

**Tech Stack:** Existing Go, zap, OpenTelemetry SDK/otelzap, Gin, and GORM dependencies.

**Spec:** The logging analysis and recommendations approved in this conversation by “按建议修改”.

## Constraints and decisions

- Retain existing `log.*` flags, output encoding and panic-level stacktrace policy; honor `disable-stacktrace`.
- Share the existing sampling policy (first 100, then every 100 within one second) across both sinks.
- Remove duplicated wrappers/interfaces/aliases and `Options.Build`; migrate all in-repository callers.
- Put request fields in standard context using a private key. Do not depend on Gin context fallback.
- Inject `*zap.Logger` into GORM. Include context in every GORM log path.
- Return startup errors rather than panicking; initialize the CLI explicitly too.
- Continue in the user's existing workspace; leave the changes available for their established separate commit step.

## Review focus

- Native zap and standard-library logs must reach the configured output.
- OTel must obey the same level and sampling decisions as local output.
- HTTP and SQL logs must carry the current span, with no cross-request leakage.
- Bad configuration must preserve the previous global logger; `New` must not change globals.
- Caller attribution and stacktrace disabling must work with native and sugared loggers.

## Tasks

- [ ] Add failing regression tests for globals, level filtering and OTLP trace linkage.
- [ ] Replace logger internals with `New`, `Init`, `L`, context fields and shared sink policy; remove redundant files.
- [ ] Migrate service/CLI startup, native call sites, HTTP context and GORM injection.
- [ ] Test HTTP/SQL context propagation, sampling, invalid configuration, caller information and concurrent initialization.
- [ ] Run focused race tests, the Go suite and diff checks; review final changes.

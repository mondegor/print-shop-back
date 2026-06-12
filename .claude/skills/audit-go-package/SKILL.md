---
name: audit-go-package
description: Multi-dimensional audit of a Go package or module — security, architecture, fault-tolerance, memory leaks, high-load/performance, and code comments. Use when asked to "audit" / "провести аудит" a package, re-audit one against a prior AUDIT.md, or do a deep correctness/robustness review beyond ordinary diff review. Dispatches one subagent per dimension, independently verifies prior "fixed" claims and rejects false positives against the source, and writes findings to <pkg>/AUDIT.md.
---

# Audit a Go package (multi-dimensional)

A repeatable methodology for a *full* audit of a Go package/module across six dimensions.
Distinct from `/code-review` (diff-focused): this audits a whole package and produces a
persistent findings document.

## When to use

- "проведи аудит пакета X" / "audit package X" / "full review of X".
- Re-auditing a package that already has an `AUDIT.md` (verify prior fixes + find new issues).
- A deep robustness/correctness pass (concurrency, leaks, panics) that ordinary diff review misses.

Do **not** use for a normal PR diff review — use `/code-review` for that.

## Dimensions (what "full" means)

1. **Security** — input validation, integer overflow, unsafe casts, secrets, injection.
2. **Architecture** — interfaces/contracts, duplication, dead code, facade/re-export hygiene,
   option/zero-value consistency, nesting/naming.
3. **Fault-tolerance** — panic recovery on every goroutine running user code; idempotent
   `Shutdown`/`Close` (`sync.Once`, no `close()` of a closed channel); context cancellation
   handling; error propagation (nothing silently swallowed).
4. **Memory leaks** — goroutine exit paths, `valueCtx` growth in loops, timers/tickers without
   `Stop`, unbounded slices/channels, batch buffers reset/reused.
5. **High-load / performance** — backpressure, channel buffering, per-item allocations in hot
   paths, lock contention, atomics vs plain fields.
6. **Code comments / docs** — stale/contradictory comments, old package names, missing godoc,
   godot (terminal period), duplicated doc text, language consistency (match the file/package),
   generated-mock headers.

## Process

1. **Read the prior `<pkg>/AUDIT.md` first** if it exists. Treat its "✅ fixed" marks as
   *claims to verify*, not facts.
2. List the package files and understand structure (`find <pkg> -name '*.go'`; read entry points,
   constructors, the main loops/goroutines).
3. **Dispatch parallel subagents — one agent per dimension group** (in a single message):
   - security + fault-tolerance
   - memory leaks + high-load/performance
   - architecture + concurrency correctness (data races; note `-race` may be unavailable without gcc)
   - comments + documentation
   Give each agent: the package path, module path, the focus checklist, the prior `AUDIT.md` path,
   and these instructions: read actual source (don't trust AUDIT.md), separate **"Verification of
   prior findings"** from **"New findings"**, cite `file:line`, assign severity, do not modify files.
4. **Personally re-verify every headline finding against the source before finalizing.** Read the
   exact lines. This catches both false "fixed" claims and agent false positives (see Lessons).
5. Consolidate into `<pkg>/AUDIT.md`: a verification table for prior findings + a new-findings list
   + a priority summary. Record any **rejected false positives** with the reasoning, so the next
   round doesn't re-raise them.

## Severity scale

- 🔴 critical — crashes the process / loses data / security hole; fix before ship.
- 🟠 major — correctness/robustness risk under realistic conditions (races, double-shutdown,
   architectural debt that will rot).
- 🟡 minor — comments, docs, consistency, low-risk edge cases.

## Go-specific checklist (per dimension)

**Fault-tolerance**
- Every goroutine running user-supplied code has `defer recover()`; recovery is per-iteration/
  per-item, not per-goroutine (a panic must not end the worker loop).
- `Shutdown`/`Stop`/`Close` is idempotent — `close(ch)` wrapped in `sync.Once` (a 2nd `close`
  panics). Remember `Shutdown` runs on the *caller's* goroutine — worker recover does NOT cover it.
- `errorHandler.Handle` / user callbacks: are *they* protected from panicking?
- `ctx.Done()` branches: is data flushed (drain + `context.WithoutCancel` + handler timeout) or lost?
- Does `Start` return `ctx.Err()` or silently `nil` (masking the cause from a supervisor)?

**Memory / high-load**
- `ctx = derive(ctx)` inside `for {}` → unbounded `valueCtx` chain. Use a loop-local var.
- `time.NewTicker`/`NewTimer` always `Stop`ped (defer). `NewTicker(d)` panics if `d <= 0` —
  is `d` floored, and is the floored value the one actually passed?
- Channels: buffered vs unbounded; blocking sends always have a `<-done`/`<-ctx.Done()` escape.
- Batch buffers bounded and reset/reused; note the dominant steady-state allocation.

**Architecture / concurrency**
- Shared mutable struct fields across goroutines — synchronized? Reason manually if `-race` can't run.
- `atomic.*` used for non-atomic read-modify-write (Load…Store) → lost updates; needs CAS loop.
- Dead packages (`grep -r '<module>/<pkg>'` returns nothing) — flag.
- Repeated contract across N types but no declared interface → only add a **public** interface
  to production code if something actually consumes it (Go favours accept-interfaces-return-structs;
  an unused exported interface is dead surface). To document/lock the shared shape, declare a
  **private** interface in each external `_test` package and assert conformance there:
  `assert.Implements(t, (*iface)(nil), &pkg.Type{})` (pattern `TestXImplementsY`). Do **not** add
  `var _ Iface = (*T)(nil)` compile-asserts in package files, and never force a new prod import
  just to satisfy an assert.
- Duplicated worker-pool / boilerplate that has already drifted → extract to `internal/...`.
- Facade re-exports public sentinel errors so callers can `errors.Is` without importing subpackages.

**Comments / docs**
- Comments contradicting code (precision claims, old names after a rename, "до секунды" vs ms).
- Missing godoc on exported symbols; duplicated doc text on option pairs (`WithX`/`WithXStrategy`).
- Generated mocks: stale `Source:`/`mockgen` header, archived deps (`github.com/golang/mock` →
  `go.uber.org/mock`).

## Verification commands

```bash
go build ./...
go vet ./<pkg>/...
go test ./<pkg>/...
golangci-lint run ./<pkg>/...     # strict project config; expect 0 issues
# go test -race ./<pkg>/...        # needs CGO_ENABLED=1 + gcc; reason manually if unavailable
grep -rn '<module>/<deadpkg>'      # confirm dead-code claims
```

## Output

Write/overwrite `<pkg>/AUDIT.md` with: header (scope, branch, date, round), "how to verify"
commands, a verification table for prior findings (claimed vs actual-by-code with `file:line`),
a new-findings list with severity, a rejected-false-positives section with reasoning, and a
priority summary. State clearly whether code was changed or the round is report-only.

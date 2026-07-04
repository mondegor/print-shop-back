---
name: go-style-guide
description: Code style guide for Go library modules following this company style. Use whenever writing, editing, or reviewing Go code in such a repo so the result matches the established conventions and passes the strict golangci-lint config. Based on the Uber Go Style Guide, adapted with grouped type blocks, English or Russian doc comments, Proto/facade patterns, table-driven _test packages.
---

# Go Library Style Guide

Conventions for Go library modules following this company style. Follows the
[Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) as a base;
the deviations and company-specific rules below take precedence. Everything here is enforced
by `.golangci.yaml` (`golangci-lint` runs strict — `make lint` must pass before done).

## Formatting & imports

- `gofumpt` with `extra-rules: true` + `gofmt` + `goimports` + `gci`. Tabs for indentation.
- Line length ≤ **160** (tab-width 4) — `lll`.
- Import groups in this exact order (`gci`/`goimports` local-prefix), blank line between:
  1. standard library
  2. third-party
  3. the module's own packages (its module-path prefix)
- No file/copyright headers (no `goheader` config). No package doc comments required
  (`staticcheck` `-ST1000`, `revive` `package-comments` disabled).
- Use `any`, never `interface{}` (`revive use-any`). Prefer `strconv` over `fmt.Sprintf`
  for simple conversions (`perfsprint`).

## Declarations — company conventions

- **Always wrap type declarations in a grouped `type ( … )` block — even a single type.**
  This is pervasive (the norm here, unlike Uber which groups only related decls):
  ```go
  type (
      // Service - описание типа ...
      Service struct {
          repo Repository
      }
  )
  ```
- Group related `var`/`const` in blocks. Predefined error catalogs are `var ( … )` blocks
  of `ErrXxx` factory protos.
- **Functional options go in their own file.** When a type uses the functional-options
  pattern (`Option` type + `WithXxx` constructors), put the option type and all its
  `WithXxx` functions in a dedicated file next to the type. Name it bare `options.go`
  **only** when the package has a single type (no ambiguity). When the package holds
  other types/files, name the file `<thing>_options.go` after the owning type (e.g.
  `session_list_options.go` next to `session_list.go`), so it's clear which type the
  options belong to — a bare `options.go` would be ambiguous there.
- **Normalize zero-value inputs to defaults in the constructor.** When a constructor takes
  config-like scalar params (durations, counts, names) with a sensible default, replace the
  zero/empty value with a package-level default **inside the constructor**, so callers can pass
  `0`/`""` to mean "use the default". Keep defaults in a private grouped `const ( defaultXxx = … )`
  block next to the type. Two shapes:
  - Plain constructor — guard each param before building the struct:
    ```go
    const defaultTimeout = 5 * time.Second

    func NewService(timeout time.Duration) *Service {
        if timeout == 0 {
            timeout = defaultTimeout
        }

        return &Service{timeout: timeout}
    }
    ```
  - Functional-options constructor — seed defaults in the initial struct literal, apply the
    options, then zero-check the rest:
    ```go
    func newOptions(opts []Option) options {
        o := options{timeout: defaultTimeout}

        for _, opt := range opts {
            opt(&o)
        }

        if o.maxAttempts < 1 {
            o.maxAttempts = defaultMaxAttempts
        }

        return o
    }
    ```
  Validation that *rejects* invalid values (e.g. an upper bound) stays separate — defaults only
  fill in zeros.
- **Flexible/self-validating types live in config; constructors take standard types.** Narrow
  YAML-bound types (`int8`, `uint16`, `uint8`, …) belong to the `wire/<comp>/config` model, where
  the type itself documents and bounds the input. Object **constructors** take the *standard* widened
  type — `int` for sizes/lengths/limits/offsets (which may legitimately be negative), `uint64` for
  identifiers — and the `wire` factory does the conversion (`int(cfg.SampleParam)`). Don't
  push narrow config types into domain signatures.
- **Default an optional interface/callback dependency to a no-op, applied *after* the options
  loop.** For an `// OPTIONAL` collaborator set via `WithXxx` (an interface or func field),
  substitute a no-op default instead of `nil`-checking it at every call site. Do the substitution
  **after** applying the options (not in the initial struct literal), so it covers both "option
  not provided" **and** `WithXxx(nil)` — the latter would overwrite a literal default and panic
  at the call site. Globals are banned, so make the no-op a tiny unexported type, not a package
  `var`:
  ```go
  type defaultAlerter struct{}

  func (defaultAlerter) SendAlert(context.Context, uuid.UUID, int) error { return nil }

  // ... in the constructor:
  for _, opt := range opts {
      opt(&o)
  }

  if o.svc.alerter == nil {
      o.svc.alerter = defaultAlerter{} // covers "not set" and WithAlerter(nil)
  }
  ```
  Then call sites stay clean: `o.svc.alerter.SendAlert(…)` with no `nil` guard.
- No global mutable state (`gochecknoglobals`) — error/sentinel `var`s are the accepted
  exception. No `init()` functions (`gochecknoinits`).
- **Repository/storage methods return simple shapes — slices or entities, never `map`.**
  A data-access method yields `([]entity.X, error)` / `([]uint32, error)`; any derived
  structure (a lookup set, an index, a grouping) is built by the **consuming layer**
  (usecase/service), not by the repository. This keeps the data-access API uniform and
  decoupled from how a caller chooses to index the rows.
- **In Postgres repositories, express "rows lacking a related row" as a `LEFT JOIN … WHERE
  x IS NULL` anti-join — not `NOT EXISTS`/`NOT IN`.** Applies to both `SELECT` and
  `DELETE … USING …` (the `USING` from-list may contain the `LEFT JOIN`). The anti-join is
  the project canon (uniform with existing queries) and keeps "find the orphan" and the
  action it feeds atomic in one statement, e.g.:
  ```sql
  DELETE FROM <sessions> s
  USING UNNEST($1::uuid[], $2::int8[]) as c(user_id, session_id)
      LEFT JOIN <auth_tokens> t
          ON  t.user_id = c.user_id
          AND t.session_id = c.session_id
          AND t.token_type = $3 AND t.token_status = $4 AND t.expires_at > NOW()
  WHERE s.user_id = c.user_id AND s.session_id = c.session_id
      AND t.auth_token IS NULL; -- anti-join: no live related row
  ```
- **Inline `LIMIT` (and similar planner-sensitive integer clauses) into the SQL text — don't
  pass them as `$N` bind parameters.** Build the clause from the concrete value
  (`mrstorage.NonZeroLimit(limit)`) and drop the arg from the
  `Query`/`Exec`/`ExecAffected`/`fetchRowsIDs` call. With `LIMIT $N` the planner can't see
  the real value at plan time and may pick a worse plan (poor row-count estimates, wrong
  scan/sort); inlining lets it plan against the actual bound. `limit` is always a
  caller-controlled `int` (a batch-size config, never user input), so this is
  injection-safe — never inline string/user-supplied values this way. `LIMIT` is normally
  the highest-numbered placeholder, so removing it needs no `$N` renumbering; keep the
  `limit` parameter (still used for the inline and for slice-capacity hints). Trade-off
  accepted: distinct limit values produce distinct SQL text (less prepared-statement reuse),
  fine because these limits are fixed configs. Canonical examples:
  `mrauth/repository/session_postgres.go`, `mrqueue/repository/queue_postgres.go`.
  ```go
  sql := `
      ... ORDER BY
          updated_at ASC
      ` + mrstorage.NonZeroLimit(limit) + `
      FOR UPDATE SKIP LOCKED ...`
  ```
- **Avoid maps in config/input DTOs — use a slice of structs with an explicit key field.**
  Конфиги и входные DTO не используют мапы. Вместо `KindLimits map[string]uint32` — слайс:
  ```go
  type (
      LimitRealm struct {
          Name       string
          KindLimits []UserKindLimit
      }

      UserKindLimit struct {
          Kind       string
          SessionMax uint32
      }
  )
  ```
- **Avoid nested maps (`map[K1]map[K2]V`) — use a flat map keyed by a small private struct.**
  Двойные мапы только по согласованию. Внутренний lookup-индекс собирается в конструкторе из
  входного слайса; ключуется приватной составной структурой:
  ```go
  type realmKindKey struct {
      realm string
      kind  string
  }

  sessionLimits map[realmKindKey]uint32
  ```
- **Name the result parameters of interface methods when the results are native types.**
  In an interface declaration, give the returns names when their types are built-in/native
  (`[]uint32`, `string`, `bool`, `int`, …) so the signature is self-documenting, e.g.
  `FetchOpenSessionIDs(ctx context.Context, userID uuid.UUID) (sessionIDs []uint32, err error)`.
  When a result is already a descriptive named type (`dto.UserScopes`, `[]entity.Session`), a
  name is optional. When other results are named, name the error `err` too.

## Comments (English or Russian godot-checked)

- Exported symbols **must** have a doc comment, in **English** or **Russian**, format
  `// Name - descript / описание.` (name, space-dash-space, then text). Doc comments
  end with a period (`godot`). Match the existing terse style.
- **Internal comments** (inside function/method bodies) may start with a lowercase
  letter; when they do, they **must not** end with a period. (`godot`'s scope is
  declarations only, so these aren't linter-enforced — follow the convention manually.)
- Document constructor params with a bulleted list when non-trivial:
  ```go
  // NewService - создаёт Service ...
  // Параметры:
  //   - repo - доступ к хранилищу данных;
  //   - handler - функция обработки результата.
  ```

## Naming

- Constructors: `NewXxx`. Receivers: short (1 letter), consistent per type
  (`s *service`, `w *wrapper`). `revive receiver-naming` enforces consistency.
- Sentinel errors prefixed `Err`, error *types* suffixed `Error` (`errname`).
  Note: error *codes* are camelCase string literals (e.g. `"errSomethingFailed"`).
- Initialisms via `revive var-naming`: `HTTP`, `JSON` (not `Http`/`Json`).
- **One canonical spelling per domain abbreviation — `2FA` for two-factor auth.** The concept
  has exactly one name: the abbreviation `2FA`. Do **not** introduce synonyms (`secondFactor`,
  `TwoFA`, `TFA`, `MFA`, `two_factor`). Because a Go package/identifier can't start with a digit,
  the spelling is context-bound but the abbreviation never changes:
  - **Exported identifiers** — treat `2FA` as an initialism, uppercase: `Auth2FA`, `User2FA`,
    `Disable2FA`, `Confirm2FA`, `Auth2FAType`, `Err2FAMustBeDisabledFirst`.
  - **Unexported locals/fields** — keep `2fa` lowercase for readability: `auth2faStorage`,
    `auth2faTableName`, `auth2faVerifier` (the compromise — uppercase mid-identifier reads
    heavy here, and `staticcheck -ST1003` is off so it's allowed).
  - **Package/dir names** (lowercase, no leading digit) — prefix with a word: `auth2fa`,
    `auth2fatype`. The verifier service package is `service/auth2fa` (not `secondfactor`).
  - **Persistence/transport** (SQL columns, JSON, YAML keys, URLs) — `2fa` / `auth_2fa` /
    `disable2fa`; never `two_fa`. Error *code* string literals follow the same form
    (`"2FAMustBeDisabledFirst"`), matching the Go var.
  The mechanism (`TOTP`) and supporting concept (`recovery`) are **not** synonyms of `2FA` —
  they keep their own names. Apply the same single-canonical-spelling rule to any future domain
  abbreviation.
- **Layer-based entity & method naming.** In the `usecase`/`service` layers, name domain
  entity values `item` / `items` and give methods **business-intent** names that describe
  the operation in domain terms. In the `repository` layer, name DB-row values
  `row` / `rows` and give methods **technical CRUD** names (`Insert`, `Update`, `Delete`,
  `Fetch`, …). The split keeps domain vocabulary in the upper layers and data-access
  vocabulary at the storage boundary.
- Import aliases lowercase `^[a-z][a-z0-9]*$`; no redundant aliases.
- `staticcheck -ST1003` is off, so some naming rules are relaxed — still follow Go idiom.

## Errors

- Follow the **Proto pattern** where used: a proto is an immutable factory built once;
  derive concrete instances via `New`/`Wrap`/`WithDetails`. Never mutate a proto after
  construction.
- Root packages act as **facades**: expose new behavior via type aliases
  (`X = subpkg.X`) and `var` function aliases (`NewX = subpkg.NewX`); implement in the
  subpackage. Prefer importing the root facade in consumers.
- Wrap errors crossing external/package boundaries — use `%w`, and compare with
  `errors.As`/`errors.Is` (`errorlint`). Boundary-wrapping itself is a **manual
  convention** — `wrapcheck` is currently disabled in `.golangci.yaml`. Never return
  `nil, nil` (`nilnil`). Don't return `nil` after a non-nil error check (`nilerr`).
- Forbidden imports: `crypto/md5`, `crypto/sha1` (`revive imports-blocklist`, `gosec`).

## Control flow (wsl_v5 / whitespace / nlreturn / revive)

- Blank line before `return`/`break`/`continue` (`nlreturn`); no leading/trailing blank
  lines in blocks (`whitespace`). `wsl_v5` governs statement cuddling — keep related
  statements together, separate unrelated ones with a blank line.
- Early return / guard clauses; avoid `else` after a returning `if`
  (`indent-error-flow`, `early-return`, `superfluous-else`, all `preserveScope`).
  Avoid deep nesting (`nestif`).
- **Keep the main positive (happy) path outside `if` blocks.** Branch on the error
  condition `if err != nil { … }` and let execution fall through to the success path —
  do **not** put the success logic inside `if err == nil { … }`. Use `err == nil` as a
  condition only as a last resort (when the negated form is genuinely not expressible or
  would obscure intent), and **always** with a comment explaining why:
  ```go
  // good — error branch guards, happy path continues unindented
  v, err := parse(s)
  if err != nil {
      return err
  }

  use(v)

  // avoid — happy path buried inside the condition
  if v, err := parse(s); err == nil {
      use(v)
  }
  ```
- Prefer `make(...)` to init maps/slices (`enforce-map-style`, `enforce-slice-style`).
  Preallocate slices with a capacity hint when length is known (`prealloc`):
  `make([]string, 0, len(x)*2)`.
- No naked returns in non-trivial funcs (`nakedret`, `bare-return`). Functions return
  ≤ 3 results (`function-result-limit`). Watch `gocyclo`/`gocritic`/`unparam`.
- No `fmt.Print*`/debug forbiddens in library code (`forbidigo`) — allowed in
  `examples/` and `_test.go`.

## Tests

- **Always external test package: `package foo_test` in a file named `foo_test.go`**
  (`testpackage`). This is a hard convention here — there are **no** internal test
  packages in the repo. Test through the **public API only**; never reach for unexported
  symbols. Do **not** use the `*_internal_test.go` / `*_export_test.go` escape hatch that
  the `testpackage` linter's default `skip-regexp` would let through — the linter allows
  it, the project does not. To cover an unexported helper, drive it through the exported
  type/constructor that uses it, not directly.
- **Always use `github.com/stretchr/testify` for assertions** (`testifylint`). Use
  `require.*` for **fatal** checks that must abort the test on failure (preconditions,
  setup, `err != nil`, non-nil results you'll dereference); use `assert.*` for
  **non-fatal** checks where the test can keep running and report further failures
  (field-by-field comparisons after a successful operation). Rule of thumb: if continuing
  the test makes no sense (or would panic) once the check fails, use `require`; otherwise
  `assert`.
- For **complex tests with mocks** where the same object/mock initialization repeats across
  most tests in the package, use a **`github.com/stretchr/testify/suite`** test suite: put
  the shared fixtures (mocks, controller, system-under-test) on the suite struct and build
  them in `SetupTest` (or `SetupSubTest`), so each test method starts from a clean,
  consistently-initialized state instead of duplicating setup. Run it via a single
  `func TestXxxSuite(t *testing.T) { suite.Run(t, new(XxxSuite)) }` entry point. Keep using
  `gomock` for the mocks themselves. For simple tests without shared setup, prefer the
  plain table-driven form below.
- For mocks use **only** `go.uber.org/mock/gomock`. Generate mocks with `mockgen`
  (`//go:generate mockgen ...`), drive expectations via `gomock.NewController(t)` and
  `EXPECT()`. Do **not** hand-write mocks or use any other mocking library.
- **Put the `//go:generate mockgen ...` directives in the `_test.go` file that consumes
  the mocks, not in the production source.** Mocks are test-only tooling, so the
  directives belong with the tests; place them right after the import block of the
  package's test file. `go generate` runs per-directory, so `-source=foo.go` (and the
  `-destination=mock/...` paths) still resolve correctly from the test file. When one
  directory has several source files generating mocks, group all their directives in the
  single package test file (e.g. `session_test.go`).
- **Always generate mocks into a nested `mock/` directory next to the consuming package**
  (`package mock`, e.g. `service/session/mock/`), one `mock/` per package that owns/consumes
  the interfaces — never into `*_mock_test.go` in the test package. The external `_test`
  package imports `<pkg>/mock` and uses `mock.NewMockXxx(ctrl)`. A mock of an unexported
  interface (generated via `mockgen -source`) is assigned to the constructor's unexported
  interface parameter structurally — the unexported type name is never written in the test.
- `t.Parallel()` at the top of every test and subtest (`tparallel`). Test helpers call
  `t.Helper()` (`thelper`).
- Table-driven with a local `type testCase struct`, named cases, `t.Run(tt.name, …)`:
  ```go
  func TestX_Method(t *testing.T) {
      t.Parallel()
      type testCase struct{ name, in, want string }
      tests := []testCase{ {name: "empty", in: "", want: ""} }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              t.Parallel()
              // ...
          })
      }
  }
  ```
- Relaxed in tests (excluded linters): `dupl`, `gosec`, `forbidigo`, `forcetypeassert`,
  `noctx`, `revive`, `unparam`.
- Benchmarks live in dedicated `*_bench_test.go` files; `go test -bench=. ./pkg/`.

## Before finishing

- Run `make lint` (or `golangci-lint run`) and `make test` (or `go test ./...`).
  The lint config is strict — treat any finding as a blocker.
- Files end with a trailing newline.

---
description: Check the project (or given paths/diff) for compliance with the Style Guide
argument-hint: "[path|package|--diff] (optional; defaults to staged+unstaged changes, else whole repo)"
allowed-tools: Bash(make lint), Bash(make test), Bash(golangci-lint run:*), Bash(go test:*), Bash(go vet:*), Bash(git diff:*), Bash(git status:*), Read, Glob, Grep, Skill
---

# Style compliance check

Audit Go code for compliance with this project's Style Guide.

## Scope

`$ARGUMENTS` controls what to check:
- empty → check uncommitted changes (`git diff` staged + unstaged); if there are none, check the whole repo
- a path or package (e.g. `path/to/pkg/`) → check just that
- `--diff` → check only `git diff` against the working tree
- `--all` → check the entire repository

## Steps

1. **Load the rules.** Invoke the `go-style-guide` skill to get the authoritative
   convention list. It is the source of truth — do not invent rules.
2. **Determine the file set** from the scope above (use `git status`/`git diff` or `Glob`).
   Ignore `examples/`, `docs/`, `vendor/`, `.github/` for source-style rules (they have
   relaxed linting), but still report obvious breakages there.
3. **Run the linters** as the objective baseline:
   - `make lint` (or `golangci-lint run` scoped to the paths if `make` is unavailable).
   - If checking behavior too, `make test` (or `go test ./...`).
   Capture and summarize findings; map each to the guideline it violates.
4. **Manual review** for things linters miss or under-enforce — go through the
   skill's checklist on the file set.

## Output

Report concisely, grouped by severity:
- **Lint failures** — copy the linter's own messages with `file:line` and the rule name.
- **Style-guide violations** — `file:line` + the specific guideline + a one-line fix.
- **OK** — note categories that passed.

Then give a short verdict: **compliant** / **needs fixes (N issues)**.

Do **not** modify any files unless I explicitly ask you to fix the findings afterward.

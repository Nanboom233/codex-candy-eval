# Independent v0.2.17 Integrated Check

## Findings fixed

- File: `scripts/test-sub2api.sh`
- Issue: The available Corepack cache could select pnpm 12 for the pnpm 9 lockfile.
- Fix: The check entry prepares and activates pnpm 9.15.9 before it runs the builder.

- The inspection coordinator could accept a start before it loaded saved settings.
  Initialization now reads settings before it publishes the coordinator.
  A failed read keeps starts blocked until a settings save succeeds.
  A focused service test covers this recovery path.

## Findings not fixed

No in-scope code defect remains from this review.
The main session owns the final check in the retained Docker fixture.
This review did not create, remove, or restart a Docker container.

The local build does not establish execution on the N150 or iStoreOS host.
The main session must report the result of its fixture smoke check separately.

## Verification

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check scripts/test-sub2api.sh check`: passed.
- Focused ESLint on `CodexEvalView.vue` and `inspection.js`: passed.
- `python3 scripts/package-sub2api.py`: passed.
- `file dist/build/sub2api-codex-eval`: reports a stripped, static x86-64 ELF executable.
- `tar -tzf dist/build/sub2api-codex-eval.tar.gz`: lists the executable, installer, and five license files.

The scoped check includes the native frontend build, `vue-tsc`, focused Go tests, touched-package `go vet`, and static Linux/amd64 build.
The JavaScript lint check declares the shared globals supplied by the existing standalone page scripts.
No additional browser run was needed for the startup-only correction.

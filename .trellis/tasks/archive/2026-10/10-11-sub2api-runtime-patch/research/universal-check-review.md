# Universal Check Review

## Findings (fixed)

- File: `install.py`
- Issue: `local_bundle()` treated any directory containing two adapter files as a complete local release. Running the repository-root installer therefore selected the checkout and failed the strict payload check instead of falling through to the tagged download.
- Fix: validate the complete 62-file allowlist and functional hash before accepting an adjacent bundle. Added a regression asserting the repository checkout is not a release while an exact package is discovered.

- File: `install.py`, `README.md`, `docs/sub2api-operations.md`
- Issue: the universal installer pointed to `v0.3.0`, while this repository already has historical `v0.3.0` through `v0.3.14` tags for the prior CPA product.
- Fix: selected unused `v0.4.0` for the universal release pin and online examples. Local tag listing and remote tag lookup showed no `v0.4.0`. This review did not create or publish a tag.

- File: `README.md`
- Issue: the developer example staged Go sources below `dist/`, so root `go vet ./...` and `go test -race ./...` discovered unbuilt package copies with an absent embedded license.
- Fix: stage under a fresh `${TMPDIR}` directory and place only the archive in `dist/`. The current universal spec also requires extracted release staging to stay outside root Go package discovery.

- File: generated review workspaces under `dist/sub2api-build/universal-20261013/`
- Issue: twenty retained generated directories were visible to root Go package discovery; the first root vet/race run failed on the `release-license-final` staging copy.
- Fix: moved each full directory, preserving all bytes, to `/tmp/codex-universal-check-20261013/`. Original paths are listed in `/tmp/codex-universal-check-20261013/moved-paths.txt`. The final repack staging is `/tmp/codex-universal-check-20261013/release-v04-current`, mapped in `final-stage-path-map.txt`. No artifact was deleted.

## Findings (not fixed)

- The required full upstream unit-race gate remains failed for both source profiles: current has 30 race warnings and 103 top-level failures; previous has 35 warnings and 126 failures. No pristine A/B run establishes attribution. Unrelated upstream fixtures were not changed.
- Previous-source native integration remains failed in three tests: the Ent export reader rejects Go export data v5, the WebSocket queue timing assertion observes less than 50 ms, and the OpenCode refresh fixture receives 429. A focused rerun of two service cases passed, but the full gate remains failed.
- Full frontend tests remain failed: current has two post-teardown `requestAnimationFrame` errors after 2,797 assertions pass; previous has six failures across two files and one unhandled table-load error. Both profiles' frontend lint/type checks passed.
- The two original user failures are represented by maintained discovery tests and the preserved pre-change test output. This review did not recreate the user's original deployment containers. The old output is test-fixture evidence, not a fresh run against those exact containers.
- The two complete 15-check Docker lifecycle runs cover Buildx and standard Docker without Buildx at the prior functional hash `5af4052c...`. The later license/package identity has separate matching-source build, install, loaded-process, status and browser evidence at `f15183c...`; the full 15-check lifecycle was not repeated after that identity change. The final installer tag and local bundle-discovery fixes were made after those lifecycle runs and do not change native runtime sources.
- Native amd64 embed cross-compilation passed, but an amd64 loaded-process lifecycle was not run. Live paid billing, production concurrency, public online bootstrap and deployed acceptance remain unverified. The v0.4.0 package is unpublished, and the required broad gates block publication.
- Compatibility is demonstrated on actual upstream revisions v0.2.14 and v0.2.15. No arbitrary future fork or release compatibility is claimed.

## Verification

- Lint: no standalone Python/Go lint target is configured. Python compilation, root vet, shell syntax and UI checks passed. The implementation evidence records frontend lint/type checks passing on both source profiles.
- TypeCheck: the native adapter type engine passed on fresh current and previous source workspaces; root `go vet ./...` also passed.
- Root `scripts/test-sub2api.sh check`: pass. This includes helper and installer suites, shell checks, root vet/race, and UI tests.
- Root `go vet ./...`: pass after generated Go workspaces were moved outside the repository.
- Root `go test -race ./...`: pass after the same move.
- Adapter `go test -race ./...`: pass.
- `scripts/test_sub2api_build.py`: 72 pass, one historical pinned-checkout skip.
- `tests/installer/test_universal.py`: 8 pass.
- `tests/installer/test_transactions.py`: 48 pass.
- Fresh current and previous source preparation each matched 73 structural edits and passed production-byte parity, including the separate UI anchor adapter.
- The adapter type engine passed against fresh current and previous workspaces, including native receiver/signature/field checks and compiled call sites.
- Focused `go test -tags=unit -race ./internal/handler -run 'TestEvalComposedNativeAdminGatewayFourFlows|TestEvalComposedCancellationAndSharedConcurrency' -count=1` passed on both fresh workspaces. It covers all four evaluation flows, exact account/group attribution, native usage/billing assertions, shared concurrency, external rejection and cancellation.
- Independent source directories and fresh semantic workspaces remain under `/tmp/codex-universal-check-20261013/source-current`, `source-previous`, `final-semantic-current` and `final-semantic-previous`; the compiled adapter helper is `native-adapter` in that directory. The first two source paths are included in `moved-paths.txt`.
- `python3 -m py_compile` and `git diff --check`: pass.

Final archive: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`.
It contains 62 regular files and is 514,343 bytes. SHA-256 is `ff096a9de4a6a7d6b69dd6f76eb107e34aab31c65549e53aa35711bed0f9c8b1`. Extraction and strict bundle validation passed; all 62 files match the fresh package staging byte-for-byte. The release tag is `v0.4.0`; the functional hash remains `f15183c55e894cd9dd1f114d3f5aabb1b364e8a33de59a24e8b2eefdbf2c6561`.

No deployment, paid request, commit, push, publication, task archival, or cleanup was performed.

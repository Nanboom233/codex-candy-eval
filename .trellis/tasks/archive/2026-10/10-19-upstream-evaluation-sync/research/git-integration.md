# Research: Git Integration and Provenance

- Query: Determine how to merge current upstream history into the native Sub2API port while preserving local history and a normal push to the fork.
- Scope: internal
- Date: 2026-10-19

## Findings

### Ref refresh at execution

The approved merge started on `sync/upstream-main-2026-10-19` from local `2aedffe`.
Its reviewed upstream parent is `7a49068`, and the fork still points to that commit.
The execution fetch found upstream `4f83676`, with two later commits:

| Commit | Treatment |
| --- | --- |
| `e120868` | Adds the all-in-one Pelican variant and schema-2 migration. Defer because the approved scope keeps Pelican unchanged. |
| `4f83676` | Releases that additional variant. Defer with its product commit. |

Use the reviewed `7a49068` parent for this integration.
The normal merge is pending without a commit. Resolve conflicts under native ownership before focused checks.
Retain the later upstream commits in the remote ref for the next synchronization.

### Commit topology and recommendation

`main` is at `2aedffe`. `upstream/main` and `origin/main` both point to `7a49068`. All three share base `6b46d4f`. Therefore, local `main` contains 21 commits absent from the fork, and the fork contains 11 commits absent locally.

Use a normal two-parent merge of `upstream/main` into a temporary integration branch created from local `main`. Do not rebase, cherry-pick the upstream range, or use an `ours` merge. Those choices either rewrite signed local history or record upstream ancestry without accurately resolving its source changes. After conflict resolution and focused checks, fast-forward local `main` to the integration result. If `origin/main` still points to `7a49068`, `git push origin main` is a normal fast-forward because the merge commit retains `origin/main` as an ancestor. Recheck the remote before pushing. If it advanced, merge its new tip first; do not force-push.

Suggested order after the plan is approved:

1. Confirm tracked state and refresh the upstream and fork refs.
2. Create `sync/upstream-main-2026-10-19` from `main`.
3. Merge the reviewed upstream tip with a normal merge commit.
4. Resolve each conflict by repository ownership, then adapt approved upstream behavior into the native patches and tests.
5. Run only checks for the changed native behavior and changed release inputs.
6. Fast-forward local `main` to the checked merge result and push `main` to `origin` without force.

The merge-tree probe reported 43 conflict entries. They include content conflicts in `AGENTS.md` and `README.md`; delete/modify conflicts for CPA source, docs, examples, and UI files removed by the native port; file-location and rename conflicts for detector data and tests; and a ModelTrace-license rename conflict. This confirms that neither blanket side selection nor an unreviewed default merge is suitable.

The native port deliberately owns production source as seven ordered patches. Keep the CPA implementation files deleted from the final repository tree. Do not restore them as a second runtime or build path. Resolve the delete/modify conflicts in `internal/plugin/`, `docs/`, `examples/`, and `scripts/test-ui.cjs` by retaining the native repository layout, then implement accepted upstream behavior in the corresponding patch and test owners. Do not apply this rule to renamed tests or licenses without reviewing their contents and current consumer.

The merge ancestry will not eliminate every future conflict. If a later upstream commit changes a CPA file that remains deleted locally, Git will report another delete/modify conflict for that path. This is an honest consequence of retaining both histories while excluding the original CPA product tree. Resolve only the incoming changed paths on each future sync; keep the adapter boundary in `patches/` rather than retaining a dormant duplicate source tree.

### Upstream commit inventory

All 11 commits after `6b46d4f` are accounted for below. Detailed algorithm, UI, and history decisions belong to the other research reports.

| Commit | Subject | Changed area and local treatment |
| --- | --- | --- |
| `0de3802` | Release v0.3.15 | Plugin release metadata and web/test files. Inspect behavior separately; do not copy the CPA release workflow. |
| `63098bd` | Replace the fingerprint test and ModelTrace with fingerprint detection | Reworks plugin Go, web, data, tests, README, docs, and example. Adapt accepted product behavior to native patch/test owners; handle detector attribution with the license owner. |
| `c363f6b` | Release v0.4.0 | Changes plugin release metadata only. Not a Sub2API server release. |
| `009c26e` | Update the README screenshot | Updates CPA documentation media. Use only if it helps the Chinese user guide. |
| `6573867` | Retry failed and unusable model answers | Changes evaluation runners, UI, tests, and README. Adapt accepted behavior into native owners. |
| `f2aeab2` | Drop the reasoning effort tooltips | Changes CPA UI copy. UI research owns whether this applies to the native page. |
| `96deb3c` | Update the README screenshot | Updates CPA documentation media again. No release or build change. |
| `d093ee8` | Version the state file and name every test consistently | Changes plugin state, migration, tests, docs, and web code. History research owns compatibility and migration decisions. |
| `dedb1bd` | Condense the repository instructions | Changes upstream `AGENTS.md`. Reconcile only instructions relevant to this repository. |
| `b4c6b60` | Send test requests as codex-tui 0.163.0 | Changes request construction, its context data, and tests. Adapt only if the native request contract supports it. |
| `7a49068` | Release v0.4.1 | Changes plugin release metadata only. Not a Sub2API server release. |

### Local commit inventory by role

All 21 local commits after `6b46d4f` are retained by the recommended merge.

| Commit(s) | Role |
| --- | --- |
| `2e01ddb`, `f84347d` | Establish the repository's Trellis specifications and platform integrations. |
| `309f2ab`, `ecf1e37`, `9ca908c`, `1da213e`, `c8d284c` | Archive completed Trellis tasks. |
| `bb86754`, `b75e23c`, `34ed91c`, `a8e10e5`, `1191b11`, `2aedffe` | Record session journals. |
| `6cac3a2` | Port the evaluator into Sub2API v0.2.17. Removes the root CPA executable, Go module, install scripts, and original `internal/plugin/` tree; adds the patch series, native tests, build entry, installer, and Sub2API release workflow. |
| `c106537` | Adjust repository agent configuration. |
| `75f5a3f` | Stabilize native account refresh and inspection controls. |
| `8c22806` | Separate saved inspection verdicts from account priority. |
| `aecf699` | Document and retain the reusable Sub2API test container. |
| `c4cc8e8` | Add confirmed-state model control and the self-extracting Go installer. |
| `47e8406` | Show saved degraded account verdicts in evaluation lists. |
| `a2c212c` | Streamline repository instructions and finish rules. |

### Release, documentation, and attribution

The upstream range changes `AGENTS.md`, `README.md`, CPA docs, examples, and plugin code. It does not change `.github/workflows/release.yml`, root `LICENSE`, or the root `licenses/` directory. Upstream tags `v0.3.15`, `v0.4.0`, and `v0.4.1` are CPA plugin releases. They do not replace the native Sub2API v0.2.17 build source or this repository's installer release flow.

Keep the native workflow's source and asset contract: it checks out `Wei-Shaw/sub2api` at `v0.2.17`, runs `scripts/build.py --source upstream --check`, and publishes `sub2api-codex-eval-installer` (`.github/workflows/release.yml:15-39,41-55`). The build applies `patches/series`, copies `tests/backend/`, and assembles the installer (`scripts/build.py:19-30,62-95`). Do not import a CPA workflow or installer.

The v0.2.17 port commit `6cac3a2` is the source-removal boundary. Its diff deletes `cmd/cpa-codex-candy-eval/`, root `go.mod` and `go.sum`, `install.sh`, `install.ps1`, the original `internal/plugin/` implementation, and old standalone docs/install tooling. It adds the native patch series and test layout. Current ownership is explicit in `.trellis/spec/backend/directory-structure.md:9-14,20-58` and `scripts/build.py:19-30`.

Review README content rather than copying it wholesale. The current Chinese README contains installer use, rollback, native account evaluation, inspection controls, and local build steps (`README.md:5-30,32-40,92-123,125-196`). Add only user-facing changes that the accepted native implementation delivers. Keep CPA-only build and standalone installation instructions out.

The attribution change is material to distribution. Upstream commit `63098bd` removes ModelTrace assets and adds `internal/plugin/data/lm-detector/LICENSE` plus detector data. This repository currently names ModelTrace in `licenses/THIRD_PARTY_NOTICES.txt:11-16`, and `scripts/build.py:71-75` embeds every root `licenses/*.txt` file in the installer. If the accepted port replaces ModelTrace with the new detector, update the notice, license file set, Chinese credits, and installer payload together. Keep ModelTrace attribution and files only if the native product still uses them. The algorithm owner must identify the exact selected dependency and license before changing those notices.

### Exclusions and focused checks

- Exclude the upstream standalone CPA launcher, `go.mod`, install scripts, and deployment flow. The native product owns a Sub2API server binary and self-extracting Go installer (`6cac3a2`; `scripts/build.py:2,55-58,79-107`).
- Exclude the upstream CPA build/release workflow. The native workflow builds from Sub2API v0.2.17 and publishes one installer asset (`.github/workflows/release.yml:15-39,41-55`).
- Exclude upstream Trellis/platform files that do not improve the native contract. Local task, source, and check ownership already lives under `.trellis/` and `tests/` (`.trellis/spec/backend/directory-structure.md:20-58`). Reconcile direct user-facing repository instructions separately.
- Do not copy CPA UI automation wholesale. The native page checks use `tests/ui/sub2api-dropdowns.cjs`; backend tests stay in `tests/backend/` and installer tests in `tests/installer/` (`.trellis/spec/backend/directory-structure.md:20-30,57-58,86-90`).
- Do not add broad build matrices or full-suite gates. Use `scripts/build.py --source <existing-sub2api-tree> --check` for affected native changes, and only the focused existing owners required by the final implementation (`.trellis/spec/backend/quality-guidelines.md:12-30,40-42`). The checks were not run during this read-only research.

## Files Found

- `.trellis/tasks/10-19-upstream-evaluation-sync/prd.md` — approved research scope, constraints, and acceptance criteria.
- `.trellis/spec/backend/directory-structure.md` — patch, test, and native source ownership.
- `.trellis/spec/backend/minimal-linux-delivery.md` — current build, installer, and published asset contract.
- `.trellis/spec/backend/quality-guidelines.md` — focused check scope and excluded gates.
- `.github/workflows/release.yml` — Sub2API source checkout, build, and installer release flow.
- `scripts/build.py` — production patch application, test placement, build, and installer assembly.
- `patches/series` — ordered native production patch list.
- `licenses/THIRD_PARTY_NOTICES.txt` — current redistributed third-party attribution.
- `README.md` — Chinese user instructions and credits.
- Git history `6b46d4f..upstream/main` — the 11 incoming commits.
- Git history `6b46d4f..main` — the 21 retained local commits.

## External References

- Upstream repository: `https://github.com/haowang02/cpa-plugin-codex-candy-eval`.
- Fork: `https://github.com/Nanboom233/codex-candy-eval`.
- Native build source: `https://github.com/Wei-Shaw/sub2api/tree/v0.2.17`.
- No external pages were needed; the compared commit objects and configured remotes supplied the required evidence.

## Related Specs

- `.trellis/spec/backend/directory-structure.md`
- `.trellis/spec/backend/minimal-linux-delivery.md`
- `.trellis/spec/backend/quality-guidelines.md`
- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/database-guidelines.md`

## Caveats / Not Found

- The merge-tree probe predicts conflicts only. It does not decide how behavior should be adapted or prove that a future merge commit builds.
- The refs were already present locally at `7a49068`; this research did not fetch, merge, change refs, stage files, build, test, or publish.
- Product decisions about detector algorithms, UI behavior, history migration, and dependency/resource cost are owned by the other research reports.

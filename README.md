# Codex Account Evaluation for Sub2API

Candy, Fingerprint, ModelTrace and Pelican evaluations for administrator account testing on Sub2API v0.2.17, Linux/amd64.

The package contains a prebuilt native server with the native sidebar and iframe view, plus the standalone `/codex-eval` page. Installation replaces `/app/sub2api` in the selected Docker container. It performs no source download or compilation.

Install over an existing v0.2.16 container to upgrade to v0.2.17, or reinstall on v0.2.17. The installer checks the candidate version before replacement.

## Install

Extract `sub2api-codex-eval.tar.gz`, enter its `sub2api-codex-eval` directory and run:

```sh
python3 install.py
```

Choose the container explicitly. Empty input selects nothing. For unattended installation:

```sh
python3 install.py install --container sub2api
```

Use Python 3.9+ and Docker CLI with permission to operate the container. The installer prints progress, backs up the original executable in a private directory under `/tmp`, replaces the binary and restarts only that container. It prints the exact backup path. Startup progress shows elapsed time and current container state/address, followed by an explicit ready message. Final failures retain request diagnostics and restore the original executable; stopped applications fail immediately. Use `--backup-dir` to choose another backup location. No `nc`, Buildx, helper image or `docker --format` is required.

Log into the existing administrator panel and choose `账户评估` beside account management in its sidebar. The native administrator view embeds `/codex-eval` in an iframe automatically; no manual menu configuration is needed. The page reads the existing same-origin `auth_token`; it has no separate key or login.

## Restore the original executable

Keep the exact backup path printed by the installer. For a running container, restore that backup with:

```sh
docker cp /path/to/backup/sub2api sub2api:/app/sub2api.restore
docker exec -u 0 sub2api sh -c 'chmod 755 /app/sub2api.restore && mv -f /app/sub2api.restore /app/sub2api'
docker restart sub2api
```

For a stopped container, copy the backup directly to `/app/sub2api`, then start it:

```sh
docker cp /path/to/backup/sub2api sub2api:/app/sub2api
docker restart sub2api
```

The backup must retain its executable permission. Check `/health` and `/setup/status` after restoration.
Default backups remain under `/tmp` until the user or system removes them.
Container recreation removes the binary replacement. Reapply the package after recreation.

## Evaluate

Choose an OpenAI group, model and accounts. Raw upstream model suggestions remain available, and you can enter a model name directly. Original prompts, sampling, scoring, renderers, options and history controls are retained. The account enable switch changes the native `schedulable` flag. Disabled accounts show the switch off and are skipped during model discovery; enabled accounts' model-fetch failures remain visible. Account errors, inactive status, rate limits and temporary restrictions display separately.

The dedicated `POST /api/v1/admin/accounts/:id/codex-eval` accepts complete streaming Responses JSON and sends to that account only. Background evaluations use the same native service directly. The selected model is sent directly, bypassing Sub2API model mappings and configuration restrictions for evaluation. Native credentials, proxy, TLS, plugin transport and Responses/Chat Completions conversion remain responsible for upstream sending. The original `/test` and ordinary user gateway remain unchanged.

Evaluations use administrator test semantics. They do not create internal users, keys or user bills, and never switch to another account. Supplier requests still consume supplier quota. Native snapshots and 401/429 handling remain; a successful evaluation does not broadly clear account restrictions.

History retains the latest 20 records per account and evaluation type. Tasks and progress stay in memory. Restart does not resume interrupted evaluations or create interrupted records. Stop controls cancel the task context and close pending responses. Clear controls retain history when saving fails. Installation does not clear history or business data.

History uses `codex-eval/history.json` under the native data directory, normally `/app/data/codex-eval/history.json` in the container.
Unreadable or null history files are protected from overwrite. Keep actual history and database files when replacing the executable.

## Automatic Candy inspection

Choose the group and account filters on the Candy tab. Inspection uses that account selection and current native memberships.
Open `设置` for inspection model, effort, maximum tests, wait range, priority, and notifications. Valid edits save automatically. There is no save button or account list on this tab.

The defaults allow three tests per account and a random 15–35 minute wait after each round. Each account gets one first test. Additional confirmation tests run only when that result indicates a possible priority change. Recovery confirmation waits 240 seconds.

Inspection uses the Candy type and plan filters and skips disabled, unavailable, and busy accounts. Its model, effort, count, and confirmation flow remain independent from manual Candy tests. Results and history appear on the Candy tab.

The Candy account header shows the next inspection time beside `账户` and `仅启用`.
Use `立即巡检` and `停止当前巡检` at the bottom of settings for manual control.

Automatic priority changes default off. Enable adjustment to show the two priority fields. Normal priority defaults to 10, with a minimum of 1. Degraded priority defaults to 99. Only accounts at these configured values can change. Lower native priority numbers are preferred.

Errors, cancellation, incomplete results, or mixed answers prevent changes. ntfy and Bark report successful changes.

Settings use the native PostgreSQL `settings` table. Empty notification secret fields retain saved values. Current progress and the next run time stay in memory. Restart reloads settings and schedules a fresh wait. It does not resume an interrupted round. Turning off automatic inspection leaves the current round running. Use the stop control to cancel it.

Container recreation removes changes in its writable layer. Reapply the package to the recreated standard v0.2.17 container.

## Development

Use an existing original Sub2API v0.2.17 source tree with frontend dependencies installed from its frozen lockfile.
Developer builds use Python 3, Git, Go, Node.js, and pnpm.
Production changes use the numbered patches listed in `patches/series`. Tests stay under `tests/`.
The build uses the upstream Go module. This repository has no separate root Go module or shared-library build.

```text
patches/                 Ordered production patches and series list
scripts/build.py         Build, check, and package entry
scripts/install.py       Repository installer source
tests/backend/           Go tests and test data, by upstream package path
tests/installer/         Installer tests and retained-container fixtures
tests/ui/                Page checks against the patched upstream tree
licenses/                License texts and third-party notices
dist/build/              Generated executable and installation archive
README.md AGENTS.md LICENSE  Project documentation and license
.github/                 Release workflow
.trellis/                Project workflow and specifications
.agents/ .codex/ .claude/ .devin/  Tool integration and skills
```

```sh
python3 scripts/build.py --source /path/to/original/sub2api --check
```

The build applies each patch in series order to a standard temporary copy of the original tree.
It uses the upstream module locks, installed frontend dependencies, and available Go caches.
It rebuilds the patched frontend, compiles the server, and creates the installation archive.
Use `--check` to attach the existing Go tests and test data to their temporary upstream packages.
That option runs the existing installer, UI, focused Go tests, and package vet checks before compilation.
Without `--check`, the build does not copy this repository's test source or test data into the temporary tree.

Fixed outputs:

- `dist/build/sub2api-codex-eval`
- `dist/build/sub2api-codex-eval.tar.gz`

The archive contains the executable, `install.py` and necessary license texts, including `LICENSE/Project-MIT.txt`.
It includes no cleanup script, development sources, tests or logs.

Local Docker checks reuse `sub2api-codex-eval-test`, including its PostgreSQL and Redis test dependencies. Run `python3 tests/installer/native_inspection_fixture.py` to load the built binary into that retained container. The command does not create another container.

See [third-party notices](licenses/THIRD_PARTY_NOTICES.txt).

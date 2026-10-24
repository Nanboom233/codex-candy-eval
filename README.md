# Codex Account Evaluation

Use Candy tests, fingerprint detection, and Pelican tests in the Sub2API administrator panel.
The plugin targets Sub2API v0.2.18.

## Install

Use a Linux/amd64 Docker host with a running Sub2API container.
Download `sub2api-codex-eval-installer` from the standalone release, then run:

```sh
chmod +x sub2api-codex-eval-installer
./sub2api-codex-eval-installer
```

Select the target container number. If necessary, select the numbered `More` entry to show all running containers.
The installer compares actual executable versions with its build version.
A mismatch or unreadable version produces a warning. Enter `y` or `yes` to continue.
Other input, empty input, EOF, or interruption cancels before backup or replacement.
Image references such as `latest` do not determine executable versions.

Keep the printed backup path. Wait for `Application ready`.
The installer restarts the selected container and waits up to 120 seconds for health and setup readiness.
Failure or interruption after replacement triggers restoration. If restoration fails, use the backup for manual recovery.
Recreate the container only when you can install the plugin again.

For multiple plugins, use the parent project's common build and `sub2api-plugins-installer` instead.

## Evaluate Accounts

1. Sign in to the administrator panel.
2. Open 账户评估 in the sidebar.
3. Select Candy, Fingerprint, or Pelican.
4. Select an OpenAI group.

Use the account button for one account. Use the list button for a batch.
A batch uses selected accounts. Without a selection, it uses all available accounts in the current filter.
Select an existing model or enter its full name.

### Candy

Select the model, reasoning effort, and attempts per account. Start the test.
The correct answer is **21**. Expand an account to read previous answers.
The success rate uses the selected model and reasoning effort.

### Fingerprint

Select a model and concurrency from 1 to 3 per account. Start the test.
Each round selects three random number challenges and compares answers with the reference models.
The result shows the closest model and ranking.
Probability appears only when all three answers are valid. Partial results show a dash.
A result that does not resemble a language-model answer has no ranking.
Expand an account to read answers, attempts, and details.

### Pelican

Select the model, reasoning effort, and original or combined version. Start the test.
The original version requests an HTML animation of a pelican riding a bicycle.
The combined version requests an SVG with Sun Wukong, Qin Shi Huang, and a polar bear.
The pelican also states the Candy answer.
The prompt view changes with the selected version. The page saves that selection.

Select a thumbnail to enlarge the result. Current results, history, and details show the test version.
Both versions share the latest 20 Pelican records per account.
You judge the image and its Candy answer.

### Records, Usage, and Cancellation

Tests consume supplier quota. A failed or unusable answer can trigger up to three request attempts.
An incorrect Candy answer does not retry. Authentication failures and rate limits do not retry.
Duration and usage include all attempts.
Each account retains the latest 20 records for each test type across page refreshes and application restarts.

The tab-bar trash button clears all test records. The account-list trash button clears records for the current test type.
Stop cancels the current test. Restart does not resume unfinished tests.
Fingerprint attempts have a three-minute limit. Candy and Pelican attempts have a ten-minute limit.

## Automatic Candy Inspection

Select a group and account filters on the Candy page. Set each account's inspection participation as needed.
Open Settings and set the inspection model and wait range. Enable automatic inspection for scheduled rounds.
Settings save automatically. The Candy page shows results, history, and the next inspection time.
Use the Settings action to start a round immediately.

### Model Control

Enable model control and select models for removal after a confirmed downgrade.
Inspection removes matching allowed models after a downgrade and restores those models after recovery.
Recovery still uses the saved removal record when model control stops or the removal selection changes.
Removed models do not block recovery inspection.

### Priorities and Notifications

Enable automatic priority changes. Set normal and downgraded priorities.
Lower numbers have higher priority.
Model control and priority changes operate independently.
Scheduled rounds and immediate inspection use these settings. Ordinary manual account and batch tests do not.
Configure ntfy or Bark in Settings to receive priority-change notifications.

### Stop Inspection

Use the Settings stop action to cancel the current round.
Disable automatic inspection to stop future scheduled rounds. This does not cancel the current round.

## Restore the Original Executable

If automatic restoration fails, use the backup path printed by the installer.
The executable backup does not include evaluation history.

The first read of older history saves `history.v1.json` under `codex-eval` in the application data directory.
An earlier history format saves `history.v0.json`.
Before a rollback, stop Sub2API and save the current `history.json` separately.
Restore the matching older history copy as `history.json`, then restore the executable as described below.
The older copy does not include records created after the upgrade.

For a running container, substitute the actual container name and backup path:

```sh
docker cp /path/to/backup/sub2api sub2api:/app/sub2api.restore
docker exec -u 0 sub2api sh -c 'chmod 755 /app/sub2api.restore && mv -f /app/sub2api.restore /app/sub2api'
docker restart sub2api
```

For a stopped container, copy directly to the executable path:

```sh
docker cp /path/to/backup/sub2api sub2api:/app/sub2api
docker restart sub2api
```

Open the administrator panel and check application readiness after restoration.

## Standalone Build

The standalone build retains its source-path interface. It needs Python 3, Git, Go, Node.js, pnpm, and original Sub2API v0.2.18 source.
Install frontend dependencies in that source first:

```sh
cd /path/to/original/sub2api/frontend
pnpm install --frozen-lockfile
```

Run this command from the plugin repository:

```sh
python3 scripts/build.py --source /path/to/original/sub2api --check
```

The existing command without `--check` remains available.
The parent project's `python3 scripts/build.py` downloads source, prepares dependencies, and always runs tests.

| Output | Use |
| --- | --- |
| `dist/build/sub2api-codex-eval` | Standalone server |
| `dist/build/sub2api-codex-eval-installer` | Standalone Docker-host installer |

Patches use lexical filename order without a series file.
Flat `tests/` contains the actual Go regression suites, UI assertions, and test support.
Native `test_*_test.go` suites declare `// sub2api-package: backend/internal/<package>` and use ordinary `testing.T` assertions.
The build stages each group in its native temporary package, preserving private-helper access.
It executes only the discovered test names, once, and stops on failure.
Python entries retain `test_*.py`. Existing standalone Go programs retain `test_*.go` with `main()`.
The UI suite executes directly as `tests/test_ui.cjs` through Node.
`SUB2API_SOURCE` supplies patched source. `SUB2API_ASSETS` supplies data resources.
Only the golden JSON remains under `assets/`. Manual container fixtures live under `tests/native_*` and do not run automatically.
The parent repository owns the shared installer transaction suite in `tests/test_installer_test.go`.
Standalone `--check` runs backend and UI suites and compiles its compatibility installer without a second transaction suite.

## Credits

- [Sub2API](https://github.com/Wei-Shaw/sub2api): the host platform.
- [CPA Codex Candy Eval](https://github.com/haowang02/cpa-plugin-codex-candy-eval): the original evaluation tool.
- [codex-candy-eval](https://github.com/haowang02/codex-candy-eval): the original Candy test.
- [LM Fingerprint Detector](https://github.com/Ikaleio/lm-detector): the number fingerprint method and reference data, under MIT.
- [LINUX DO](https://linux.do/): community support.

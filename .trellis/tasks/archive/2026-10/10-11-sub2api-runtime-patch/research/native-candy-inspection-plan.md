# Native Candy Inspection Plan

This document records the remaining inspection work within the complete plugin migration to Sub2API. The user approved implementation with `批准实施` after the integrated plan. Earlier decisions and execution reports remain valid. See `migration-history-and-decisions.md` for the complete history.

## Source behavior

The original `examples/inspection.py` is a one-shot Candy inspection script. An external cron decides when to invoke it. Each invocation performs one first test per eligible account. It performs up to two additional tests only when the first result could change the configured normal or degraded priority. Recovery confirmation waits 240 seconds. Three valid matching results are required for a priority change. The script skips disabled, unavailable, busy, and optionally Free accounts; it matches configured include/exclude patterns against name and email. The script polls every 15 seconds and waits at most 30 minutes for a test phase.

The original script does not schedule fingerprint, ModelTrace, or Pelican. The original plugin inspection HTTP endpoints support more test types, but they are not part of the scheduled script behavior.

## Native decision

Move the scheduled behavior into the existing Sub2API process. Docker is the deployment boundary, so the coordinator runs with the application and uses native services. Do not require host cron, a second container, Python, Docker socket access, or a persistent admin credential.

Add a `设置` tab to `/codex-eval`, styled like native `frontend/src/views/admin/SettingsView.vue`. Save inspection settings through the native `SettingService` and `SettingRepository` in the existing PostgreSQL `settings` table. This replaces the proposed `inspection.json` design. Keep the next run time in memory.

Settings:

- automatic inspection switch;
- literal model and effort;
- maximum tests per account, meaning first test plus at most `N-1` confirmations;
- minimum and maximum random wait minutes;
- include and exclude account patterns;
- skip Free switch;
- automatic priority update switch;
- normal and degraded priority values;
- ntfy and Bark fields.

Defaults are disabled, 3 maximum tests, and 15/35 minute waits. Equal wait bounds are valid and create a fixed delay. Each completed scheduled run draws a new delay. A manual immediate run starts now and schedules the next delay after it ends.

## Decisions

- Candy only for automatic inspection.
- Lower Sub2API priority numbers are preferred; do not copy CPA's larger-is-preferred values.
- Only accounts currently at one of the two configured priority values are eligible for automatic updates.
- Automatic updates default off and recommendations remain visible when off.
- Disabled/schedulable-false accounts do not trigger model discovery errors.
- Cancellation, timeout, incomplete or inconsistent results never update priority.
- Notification failures are reported separately and do not rerun or undo a successful update.
- The process restores settings after restart but does not resume interrupted work or synthesize history.
- Use native administrator settings routes and the existing key/value repository. No new table or storage service is needed.
- Follow native secret-field behavior. Return configured status instead of saved credentials. Empty secret input preserves the saved value.
- Response-model comparison is outside the plugin contract. Its absence is not a defect or a required test.
- No new administrator authentication mechanism, external script, cron export, cleanup tool, or inspection history database.

## Native settings evidence

The original source defines `SettingRepository` in `backend/internal/service/setting_service.go:99-106`. Its `GetMultiple` and `SetMultiple` methods use the existing database table through Ent (`backend/internal/repository/setting_repo.go:85-115`). The schema names this table `settings` (`backend/ent/schema/setting.go:28`). No schema migration is required for additional setting keys.

Native SMTP settings preserve an existing password when the update is empty (`backend/internal/service/setting_update.go:196-202`). Read responses return `SMTPPasswordConfigured` (`backend/internal/service/setting_parse.go:338-344`). Apply this existing pattern to inspection notification credentials. Native administrator settings routes are registered in `backend/internal/server/routes/admin.go:567-571`.

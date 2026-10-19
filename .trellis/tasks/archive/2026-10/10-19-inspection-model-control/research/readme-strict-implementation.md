# Strict README rewrite

## Result

The complete README now follows the user's strict `output` rules.
The rewrite covers installation, administrator access, evaluations, inspection, restoration, and development.
Numbered procedures separate reader actions. Short paragraphs describe automatic behavior.
The rewrite removes repeated container-recreation advice and separates manual restoration commands.

The document preserves versions, commands, paths, defaults, Chinese labels, supplier quota, history, authentication, and recovery conditions.
It states that scheduled inspection and `立即巡检` apply model control and optional priority adjustment.
Ordinary manual single-account and batch evaluations do not apply these controls.
Saved account removals control recovery after switch-off or global-list changes.
Recovery inspection remains available after model removal.

## Sources

| Source | Use |
| --- | --- |
| `/Users/minami/.codex/skills/output/SKILL.md` | Strict language and procedure rules |
| `README.md` and its staged version | Complete original facts and commands |
| `implement.jsonl`, `prd.md`, `design.md`, and `implement.md` | Current requirements and approved rewrite scope |
| `.trellis/spec/backend/runtime-installation.md` | Administrator session, literal model requests, evaluation behavior, and history |
| `.trellis/spec/backend/native-inspection.md` | Inspection defaults, confirmations, controls, and recovery |
| `.trellis/spec/backend/database-guidelines.md` | History, settings, and saved removal data |
| `.trellis/spec/backend/minimal-linux-delivery.md` | Installer, backups, restoration, and retained fixture |
| `.trellis/spec/backend/directory-structure.md` | Patch, test, and output locations |
| `.trellis/spec/backend/quality-guidelines.md` | Documentation-only verification scope |
| `research/model-storage.md` | Whitelist entries, aliases, unrestricted accounts, and passthrough |
| `research/go-installer-parity.md` | Installer transaction and fixture launcher boundary |
| `scripts/installer/main.go` | Numbered selection, immutable IDs, versions, readiness, and restoration |
| `scripts/build.py` | Build command, checks, temporary modules, outputs, and embedded licenses |
| `.github/workflows/release.yml` | Single installer release asset |

## Documentation checks

A read-only sentence check found no prose sentence above 20 words.
The check excludes code, paths, commands, and other literals from word counts.
No paragraph exceeds six sentences. Prose contains no semicolons or contractions.
The passive-voice and auxiliary-pattern scan found no matches.
Manual review checked meaning, active voice, verb choice, consistent terms, procedure structure, and conditions before commands.

Every original shell command remains in the rewrite.
The local Markdown link and all listed repository paths exist.
The license source contains the five texts that the build embeds.
`git diff --check -- README.md` passes.

## Handoff

Only `README.md` and this record changed under this assignment.
The implementation leaves staging and commits to the main session.
No product build, runtime check, Docker operation, or product test ran for this documentation change.
The candidate is ready for independent documentation review.
No material source issue remains.

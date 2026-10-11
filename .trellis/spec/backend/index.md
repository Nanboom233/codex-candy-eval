# Plugin Development Guidelines

> Source-backed rules for the CPA evaluation plugin and its embedded management page.

## Scope

This is a single Go module with a C ABI bridge and embedded plain JavaScript, HTML, and CSS. It has no database, ORM, standalone web server, or frontend framework build.

Use this spec layer for plugin handlers, model requests, persisted history, concurrency, and backend-facing UI contracts. Read the actual browser modules before changing UI behavior.

Repository-wide rules remain in [AGENTS.md](../../../AGENTS.md). Release checks and packaging remain in [release.yml](../../../.github/workflows/release.yml).

## Guidelines Index

| Guide | Purpose | Status |
| --- | --- | --- |
| [Directory Structure](./directory-structure.md) | Package ownership, file placement, C bridge, and embedded assets | Source-backed |
| [JSON Persistence](./database-guidelines.md) | Snapshot compatibility, lock ownership, history limits, and file safety | Source-backed |
| [Error Handling](./error-handling.md) | Host envelopes, management errors, request failures, and retry boundaries | Source-backed |
| [Quality Guidelines](./quality-guidelines.md) | Concurrency, test isolation, ABI ownership, UI checks, and release gates | Source-backed |
| [Diagnostics and Logging](./logging-guidelines.md) | User-visible diagnostics and sensitive-data exclusions | Source-backed |

The persistence guide keeps its original filename for link compatibility. Its content describes JSON files, not SQL or an ORM.

## Pre-Development Checklist

1. Read the active task's requirements and approved planning artifacts.
2. Read `AGENTS.md`.
3. Identify the smallest change boundary.
4. Read the directory guide and the specific guides for the affected behavior.
5. Inspect the owner files and adjacent tests before editing.
6. Search for existing helpers before adding scheduling, response, history, or model-execution logic.
7. Read the [shared thinking guides](../guides/index.md) when changes span boundaries or repeat existing behavior.
8. Check host and UI consumers before changing management or C ABI contracts.
9. Preserve unrelated work.
10. Keep private deployment values out of files and reports.

The shared guides contain upstream Trellis examples. Apply their general reasoning to this plugin. Do not assume their upstream package paths exist here.

## Quality Check

For changed Go files, run `gofmt`, then:

```sh
go vet ./...
go test -race ./...
```

For embedded UI changes, also run:

```sh
node scripts/test-ui.cjs
```

Review JSON tags, response shapes, lock ownership, worker shutdown, and persistence failure behavior for the affected flow.

Default Go checks do not build the `cshared` bridge. For ABI or build changes, follow the tagged shared-library build instructions in `AGENTS.md`.

Plugin code, UI, build, and installer changes require deployment and authenticated end-to-end checks. Documentation and local setup changes are exempt.

For spec-only changes, verify source claims, local links, and the absence of unfinished content. Report incomplete checks instead of treating them as successful.

## Documentation Rules

Write project specs in English. Keep assistant replies in Chinese. Use `/output` style: short sentences, active voice, common words, and one name for each concept.

Limit instructions to 20 words. Limit descriptions to 25 words. Use one instruction per sentence. Do not use semicolons or contractions. Keep each paragraph to one topic and at most six sentences.

Describe current repository patterns. Cite real owner files or tests. Do not invent dependencies, frameworks, migration systems, or production guarantees.

For commits, follow the [combined Trellis and git-commit process](./quality-guidelines.md#commit-and-finish-work). Record spec changes before the work commit.

# Trellis Orchestration Correction

Date: 2026-10-16.

## Findings

The main session dispatched three implementation agents before it completed the full workflow and spec review. The context manifests contained research and task artifacts but no actual spec entries. This did not satisfy workflow Phase 1.3.

All three workers stopped after provider 429 or 503 errors. No inspection source files appeared before the correction. A proposed API message is design input, not implementation evidence.

## Correction

The main session read the complete workflow, implement/check/research profiles, relevant skills, active backend specs, and shared guides. Both manifests now contain actual applicable specs and research. Task artifacts load separately through the agent context protocol.

One implement agent owns the remaining settings, runtime, page, and delivery work. An independent check agent follows it. The main session owns spec synchronization, work-commit review, and task closure.

The integrated plan remains approved. This correction changes orchestration and context, not the requested functions or acceptance scope. Existing baseline and passed checks remain reusable.

## Applicable Contracts

- `.trellis/workflow.md`: Phase 1.3, 2.1, 2.2, 3.3, and 3.4.
- `.codex/agents/trellis-implement.toml`: context fallback and recursion guard.
- `.codex/agents/trellis-check.toml`: independent review and local fixes.
- `.trellis/spec/backend/index.md`: active spec routing and check scope.
- `.agents/skills/trellis-update-spec/SKILL.md`: seven sections for cross-layer code-spec contracts.

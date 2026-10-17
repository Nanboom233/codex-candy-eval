# Docker Formatting Compatibility Fix

User reports: unknown flag: --format, followed by See docker --help.
The live product has formatted container enumeration in both install.py and cleanup.py, and formatted Buildx enumeration in cleanup.py.
The exact remote emitting command is not established by the short output. Fix the shared CLI compatibility class without asking for another dump.

Use Docker quiet-ID container listing instead of Go-template --format, then existing JSON inspect for facts.
Avoid --format in Buildx inventory. Detect optional Buildx availability without producing an expected unsupported-flag failure.
If Buildx is available but only table output is supported, parse conservative exact builder names and retain ambiguous rows.
Missing optional Buildx must retain builders while allowing ordinary container/image/temp cleanup.
Daemon, permission, invalid inspect and unexpected command failures remain errors with complete original output and tracebacks.
Do not weaken exact resource selection, refreshed references, immutable container identity or business-data protection.
Do not implement removed source/build/binary-attestation paths.

Reproduce the exact unsupported-format message through a real subprocess boundary.
Cover both installer enumeration and cleanup inventory, absent Buildx, present old Buildx lacking format support,
and unrelated Docker failures. Use meaningful old-CLI fake processes plus current real Docker where available.
Do not describe subprocess simulations as real iStoreOS acceptance.
Repackage fixed dist/build/sub2api-codex-eval.tar.gz using existing dist/build/sub2api-codex-eval unchanged.
Verify actual archived entry points, raw eight-file/source/license/binary parity and affected regressions.
No runtime compilation, broad native/frontend checks, production deployment, paid request, commit or publication.
Original implementer fixes/tests/packages, then original checker reviews. Main owns specs/task/final artifact verification.

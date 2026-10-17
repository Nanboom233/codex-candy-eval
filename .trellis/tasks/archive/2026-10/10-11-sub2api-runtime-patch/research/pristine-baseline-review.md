# Pristine baseline evidence review

Main reviewed the completed execution report and independently checked evidence on 2026-10-14.
Task remains planning. No product implementation or package acceptance follows from this review.

## Independently verified

- Both source trees match every one of the official archive's 4,322 original files byte-for-byte; no original file is missing or modified.
- The archive SHA256 matches the verified source contract. All 165 frozen maintained product files and both canonical delivery artifact digests match the preservation manifest.
- The pure artifact SHA256 is `9bc90c95a5d8811e1e93aff55175ab0be7d999740fd4e5cb632d48d586a5a2bf`. `file` and `go version -m` confirm static x86-64 ELF, Go1.27.2, Linux/amd64, CGO0, embed and trimpath. The actual loaded version log returns 0.2.16 and the verified full commit.
- All reported gate command/environment/exit/duration/log records match their original run records. Both complete stdout and stderr files exist; their sizes and hashes are preserved in `pristine-baseline-main-review.json`.
- Independent parsing of raw Linux unit stdout yields 11,806 top-level passes and 18 skips, with zero failures; including subtests, 22,907 passes. Integration yields 7,127 top-level passes and 15 skips; including subtests, 14,061 passes and 16 skips. Counts match the report.
- Independent parsing of original race JSON yields 140 top-level failures, 149 including subtests, and 28 race warnings. Both focused reproductions retain full race frames and nonzero exits.
- Smoke records cover the actual compiled binary's setup status, configured health/status and embedded frontend responses. All five HTTP outcomes are 200. They are emulated Linux/amd64 process checks on an ARM64 Docker VM.

## Result limits retained

The compiled baseline and frontend checks pass. Required unit coverage passed in a Linux ARM64 Alpine environment with a writable byte-identical scratch tree, explicit DNS-host prerequisites and 56 cached passing package results. It is not an Ubuntu CI result or a native Linux/amd64 test run. Native macOS full unit attempts retain DNS failures.

Integration's successful exit does not cover its explicit skipped tests. The optional plugin-runtime test requires a package and remains skipped in the pristine baseline. No provider credentials or paid requests were supplied.

Supplemental full race remains failed. The pristine race reproductions establish failures without local patches, but do not individually attribute or excuse the historical patched failures. No claim that all tests pass is justified.

The checker independently reported source/product/lock parity and skip/platform limits during its preliminary review. Its final turn failed twice with model transport HTTP503; main performed this final evidence review instead. This transport failure is separate from all baseline test results.

Complete raw diagnostics remain in the retained system-temp workroot linked by the execution report. No source fixes, assertion changes, race suppression, package promotion or task activation occurred. Baseline evidence collection is complete; subsequent remediation/integration choices remain planning work.

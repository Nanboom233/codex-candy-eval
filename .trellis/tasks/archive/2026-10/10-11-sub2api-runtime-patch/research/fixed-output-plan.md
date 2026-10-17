# Fixed Delivery Output Revision

The user requests one fixed output directory and stable artifact names.
Canonical archive: dist/build/sub2api-codex-eval.tar.gz.
Canonical raw patched binary: dist/build/sub2api-codex-eval.
Keep the existing eight-file archive payload and v0.2.16/amd64 runtime unchanged.
The inner binary name may remain sub2api-v0.2.16-linux-amd64.

Implementation owns installer asset URL, packager defaults, CI output/upload/publish paths, current docs and affected tests.
Copy the already reviewed patched binary to the canonical raw path. Repackage changed installer bytes with the maintained explicit Python archive writer.
Do not compile again. Use standard temporary staging where practical, or one stable staging directory under dist/build.
Remove dated/version-specific paths from active delivery commands. Keep historical archives and reports as evidence.
Online asset basename must agree with the fixed canonical archive basename.
No business runtime, cleanup semantics, production deployment, paid request, commit or publication changes.

Original implementer executes. Original checker verifies changed code/current references, actual canonical archive, original binary identity,
raw eight-file allowlist, installer extraction, archived entry points and appropriate affected tests.
No new complete native build/race/frontend runs are needed for output/asset-path changes.
Retain previously reported complete native race failure and external limits.
Main synchronizes specs/task evidence and delivers the fixed archive link after review.

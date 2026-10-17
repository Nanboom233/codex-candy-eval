# Deployment Evidence Reader Specification Review

Date: 2026-10-12. Scope: final specification review of patch 005.

## Context and Verdict

Context mode: native Trellis hook injection. The visible hook output was
truncated. The saved hook artifact was read in bounded sections, and the
registered reader plan, updated plan, and final delivery were read directly.

Verdict: PASS after the accepted unused-runner correction. No product files
were changed by this reviewer. The distinct quality review may proceed.

## Finding Resolved

The first review found that the reader rejected every declared runner not
referenced by a case. The updated frozen plan now explicitly permits valid,
unique unused runner descriptors. The corrected reader validates each
descriptor, permits unused ones, and still requires every case destination to
reference a declared runner. Exact external case coverage remains enforced.

The new `TestEvalDeploymentEvidenceAllowsValidUnusedRunner` test accepts a
valid unused descriptor and rejects it when its required client version is
empty. No findings remain unfixed for this reader scope.

## Final Delivery Identity

- Patch SHA-256: `3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb`.
- Fresh replay: `dist/sub2api-build/scratch-deployment-reader-specfix-1wk5jtxn`.
- The previous hash `991d323e24d381c75b1754a84d09b55a791393f30e039b20713d693fd90f4dcc` and replay `scratch-deployment-reader-final-vnhun2qq` are superseded.

## Final Verification

These checks ran on the corrected maintained generated backend:

- `gofmt -l` on the reader and adjacent test returned no files.
- `go test -tags unit -race ./internal/service -run '^TestEvalDeployment' -count=1` passed.
- `go test -race ./internal/service -run '^TestEvalDeployment' -count=1` passed compilation and reported no matching tests.
- `go vet ./internal/service` and `go vet -tags unit ./internal/service` passed.
- Corrected reader source and test matched the fresh replay byte-for-byte.
- Full maintained/replay backend and frontend tree comparisons passed. The replay evidence records 3,258 backend files, 931 frontend files, and 28 core files.
- The maintained and replay core copies matched `internal/eval`. Backend `go.mod`, `go.sum`, and frontend `pnpm-lock.yaml` matched.
- Frozen patch hashes 001-004 matched their registered values. Patch 005 matched the final hash above, and its patch contains only the reader and adjacent test.

Root vet/race and copied-core vet/race passed during the first pass, before the
scoped service-only correction. They were not rerun because the correction
does not touch those packages. The full generated Sub2API service race suite
was not rerun; the task records reproduced pristine upstream fixture race
failures and the user declined unrelated fixture repair. No full-suite pass
is claimed.

## Quality Handoff

The quality reviewer should independently replay the final patch identity
above and review the corrected artifact. The later producer must encode all
four case tuple members explicitly, including empty values. The trusted
current-facts provider must supply deployment, source, patch, image, binary,
and ingress facts independently. The coordinator must reload evidence at its
specified lifecycle checks. Patch 005 adds no producer, provider, coordinator,
writer, network call, or readiness change.

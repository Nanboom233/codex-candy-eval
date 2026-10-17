# Evaluation core and build/install scope inventory

Read-only planning audit, 2026-10-14. **65 pending items: C01–C34 and T01–T31.** Each item requires a new individual retain/remove decision; the recommendations below authorize no edit. Main owns approved race-gate removal and the task/spec updates. The parallel native inventory owns individual rules/runtime templates.

The actual build path is `prebuilt_dev.py` → explicit `SOURCES` → `adapter/build.py` → structural `rules.json` / `ui-rules.json`, `adapter/runtime` and copied `internal/eval`. The eight-file delivered payload contains the ready binary, `install.py`, `cleanup.py` and five license texts. It contains no core/adapter source or tests. Current `install.py` performs no source resolution or compilation.

Historical top-level and `production/` v0.2.15 patches are not current runtime inputs. The v0.2.15-named **test-overlay/** is still an active developer test input. Most `sub2api_build.py` source/patch/workspace machinery is unused, but its checksum/facts helpers remain imported by current packaging/CI. `adapter/Dockerfile` is not executed by the current prebuilt path but remains a functional-hash input. These dependencies rule out deleting files merely because their names look old.

Historical PRD references below are secondary records that attribute requirements to the user; they are **not original user transcripts or current approval**. Specific unrequested mechanics are marked absent. The original CPA source provides behavior comparison, not permission to preserve every copied feature.

Read-only evidence: numbered source reads, `diff` and `rg` reference search. Data trees, Fingerprint browser code and all vendor renderers compare byte-identical; both pure algorithm files differ only the package declaration. Root UI harness targets original CPA assets. No product edit/build/test, target access, resource mutation, cleanup, deployment or publication occurred. Existing canonical artifacts remain untouched. This inventory does not establish the installed evaluation's failure cause.

## Requirement provenance keys

Historical PRD records are secondary requirement attribution, not original user transcripts or current approval. Each item references a key below. Its JSON record expands the same provenance and complete source paths.

- **preserve**: Historical PRD R1–R2 attributes the request to port all four tests and preserve their semantics. No original user transcript is contained in that artifact; this is requirement provenance, not current approval. [prd.md:23–24](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L23).
- **native**: Historical PRD records native billing/concurrency, exact group/account routing and native administrator auth; its latest delivery section preserves those behaviors. No current individual retain decision exists. [prd.md:25–29](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L25); [prd.md:223–228](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L223).
- **history**: Historical R8 records latest-20 history, selected group and refresh/restart persistence; the later delivery record mentions cancellation/history. It does not establish user approval of every durability mechanism. [prd.md:30](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L30); [prd.md:228](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L228).
- **ui**: Historical R13/R16 records preservation of four tabs/controls, native theme and first activation using existing web login. Exact implementation details remain unapproved individually. [prd.md:35](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L35); [prd.md:38](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L38).
- **none**: No direct user requirement authorizing this specific extra mechanism was found in the inspected task artifacts. An earlier spec or implementation is insufficient retention approval.
- **prebuilt**: Historical latest-delivery section explicitly attributes v0.2.16 Linux/amd64 binary replacement, /app/data layout, no on-host compilation and minimal payload to the user. This remains past requirement evidence, not current retention approval. [prd.md:223–231](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L223).
- **selection**: Historical PRD attributes explicit container selection, single Python entry and concise complete raw diagnostics; later recorded user failure requires Docker formatting compatibility. Exact helper/validation mechanisms remain pending. [prd.md:38–39](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L38); [prd.md:241–243](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L241).
- **transaction**: Historical PRD records reversible container-local replacement, original backups, restart/recovery and install/status/rollback/uninstall. No direct transcript or individual approval of every phase/check is available here. [prd.md:31–32](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L31); [prd.md:227–228](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L227).
- **cleanup**: Historical latest delivery records an explicit independent cleanup.py for old build leftovers, preservation of business resources/probe helper and no global prune. Earlier artifact cleanup approval explicitly retained maintained test sources. [prd.md:234](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L234); [prd.md:170](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L170).
- **output**: Historical latest-output record explicitly attributes fixed dist/build binary/archive paths and reuse without recompilation to the user. [prd.md:238–239](../../../../.trellis/tasks/10-11-sub2api-runtime-patch/history/before-pristine-baseline-planning-2026-10-14/prd.md#L238).

## C01 — Host-neutral evaluation API and CPA separation

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/eval.go:13–57](../../../../internal/eval/eval.go#L13); [internal/eval/model.go:20–30](../../../../internal/eval/model.go#L20); [internal/plugin/model.go:90–110](../../../../internal/plugin/model.go#L90).

**Behavior:** Introduces positive numeric group/account targets, display-only account metadata and a ModelExecutor interface. CPA credential enumeration, host callbacks, C ABI and inspection scheduler are absent from the copied core.

**CPA comparison:** Original plugin used auth.ID/provider and host.model.execute_stream; internal/plugin remains separate legacy source.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native runtime executor and resolver templates; all four core evaluators; T17 source allowlist.

**Deletion effect:** Deleting the interface/types breaks all native evaluators and runtime imports. Returning to original CPA executeModel requires unavailable CPA host callbacks.

**Recommendation:** Retain the required host boundary, subject to individual approval.

## C02 — Requested/effective identity and outcome fields

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/eval.go:59–79](../../../../internal/eval/eval.go#L59); [internal/eval/eval.go:133–136](../../../../internal/eval/eval.go#L133); [internal/eval/runs.go:351–370](../../../../internal/eval/runs.go#L351).

**Behavior:** Adds group/account/run identity, display name, requested and effective model/effort, and status to every result; synthetic failures carry the same metadata.

**CPA comparison:** CPA result structs stored model/effort and credential-keyed history, without ResultInfo.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** C03 model normalization; C19 history schema; UI state adapters.

**Deletion effect:** Removing fields requires coordinated native/manager/history/UI consumers; historical account-only result maps cannot distinguish selected groups without the result field.

**Recommendation:** Retain routing/result identity; review requested/effective duplication separately if simplification is desired.

## C03 — Model suffix normalization and supported effort changes

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/eval.go:223–267](../../../../internal/eval/eval.go#L223); [internal/eval/web/app.js:6–9](../../../../internal/eval/web/app.js#L6); [internal/plugin/web/app.js:6–9](../../../../internal/plugin/web/app.js#L6); [internal/plugin/model.go:24–27](../../../../internal/plugin/model.go#L24).

**Behavior:** Parses known parenthesized suffixes; rejects parentheses otherwise, unknown suffixes and explicit conflicts; accepts none/minimal/low/medium/high/xhigh and omits none on wire. UI removes original max choice while keeping its five other choices.

**CPA comparison:** Original validModelName allowed arbitrary parentheses when requested; CPA owned suffix translation and UI offered max.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C02 effective metadata; C04 validation; GatewayExecutor body checks; state validator.

**Deletion effect:** Deleting normalization leaves model suffixes on native wire and breaks validation callers; relaxing/removing max rejection changes selection semantics.

**Recommendation:** Retain plain-wire normalization; submit the max-removal/restricted effort list as a concrete behavior choice, since no specific user request for dropping max was found.

## C04 — Evaluation-specific option validation

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/eval.go:95–132](../../../../internal/eval/eval.go#L95); [internal/eval/runs.go:141–199](../../../../internal/eval/runs.go#L141).

**Behavior:** Validates type, modes and sample concurrency both for single and managed runs. Fingerprint forces low and rejects explicit none/non-low suffixes. ModelTrace rejects nondefault effort; only Candy accepts batches. Candy/Pelican reject negative concurrency in managed admission.

**CPA comparison:** CPA handlers validated modes and concurrency, Fingerprint had no effort field; native adds shared validation and explicit effort restrictions.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C03; Fingerprint/ModelTrace collectors; C18 managed batching.

**Deletion effect:** Removing constraints can invoke invalid mode slices or run unsupported batch/effort combinations; duplicate validation can be reviewed without deleting behavior.

**Recommendation:** Retain inherited mode/concurrency rules; explicit suffix conflict/none rules and duplicate entry validation are implementation additions requiring individual review.

## C05 — Canonical in-process Responses HTTP executor

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/http.go:21–70](../../../../internal/eval/http.go#L21).

**Behavior:** Wraps runtime-owned RoundTripper, POSTs to fixed codex-eval.internal/v1/responses, checks target/body/model/stream/effort consistency and passes only fixed Codex headers. No credential or directed target HTTP header is added; redirects are rejected.

**CPA comparison:** Replaces CPA host bridge and earlier loopback-token executor; there is no current NewHTTPExecutor constructor.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native runtime transport binds target via private context capability; C03; C07 SSE reader.

**Deletion effect:** Deleting the executor breaks native dispatch transport; removing individual redundant body/header checks changes admission but does not itself remove native authorization.

**Recommendation:** Retain the canonical executor; individual validation checks have no separate user-level authorization beyond native dispatch.

## C06 — Safe errors, panic recovery and cancellation answer suppression

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/model.go:32–43](../../../../internal/eval/model.go#L32); [internal/eval/model.go:136–149](../../../../internal/eval/model.go#L136); [internal/eval/http.go:70–100](../../../../internal/eval/http.go#L70).

**Behavior:** Returns status/code errors without upstream bodies/URLs; keeps parsed error-response usage, closes body on cancellation/success and suppresses answers after cancellation. Executor panic becomes retryable status-zero executor_panic; String/GoString hide internals.

**CPA comparison:** CPA recovered panic with raw panic value and exposed detailed host/SSE bodies.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Retry logic C10; SSE usage C09; result errors; native runtime.

**Deletion effect:** Removing safe error codes can expose provider/private data; removing panic-to-status-zero changes retries. Deleting cosmetic String methods only changes diagnostic formatting.

**Recommendation:** Retain body ownership and safe errors; review panic retry behavior and formatting helpers individually if removing extras.

## C07 — Incremental SSE terminal and item handling

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/sse.go:52–151](../../../../internal/eval/sse.go#L52); [internal/eval/sse.go:167–215](../../../../internal/eval/sse.go#L167); [internal/plugin/model.go:186–249](../../../../internal/plugin/model.go#L186).

**Behavior:** Dispatches complete blank-line SSE events, supports multiline data/event-name fallback, requires response.completed, rejects malformed/error/incomplete/failed events and DONE-without-terminal. Reconstructs empty terminal output from completed items, replacing equal item IDs; excludes tool arguments and rejects empty/tool-only answers.

**CPA comparison:** Original buffered entire host stream, ignored malformed JSON lines, appended done items and then checked terminal status.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** GatewayExecutor; responseItem/responseBody; four evaluators.

**Deletion effect:** Deleting parser prevents Responses output consumption. Relaxing malformed-event/blank-line checks could accept partial streams; item-ID replacement specifically prevents duplicated text.

**Recommendation:** Retain terminal/text semantics inherited from CPA; explicitly review stricter malformed-event rejection and ID replacement as additions.

## C08 — SSE and request/header size bounds

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/http.go:12–19](../../../../internal/eval/http.go#L12); [internal/eval/http.go:35–36](../../../../internal/eval/http.go#L35); [internal/eval/http.go:61–63](../../../../internal/eval/http.go#L61); [internal/eval/sse.go:65–68](../../../../internal/eval/sse.go#L65); [internal/eval/sse.go:96–100](../../../../internal/eval/sse.go#L96); [internal/eval/sse.go:134–135](../../../../internal/eval/sse.go#L134); [internal/eval/sse.go:158–162](../../../../internal/eval/sse.go#L158); [internal/eval/sse.go:205–206](../../../../internal/eval/sse.go#L205).

**Behavior:** Adds 1 MiB request/line/event bounds, 64 KiB per allowed header value, 4 MiB item/output accumulation, 16 MiB total stream and 65,536 dispatched events.

**CPA comparison:** CPA response reader had no equivalent cumulative/event/line bounds.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** C05 executor and C07 parser; observed usage retained on bound failures.

**Deletion effect:** Deleting bounds removes new rejection cases and allows unbounded buffers; deleting just one constant breaks its references.

**Recommendation:** Review each bound with the user; recommend keeping bounded streaming but do not infer approval of these exact values.

## C09 — Observed-only token accounting through failed attempts

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/sse.go:20–50](../../../../internal/eval/sse.go#L20); [internal/eval/sse.go:81–84](../../../../internal/eval/sse.go#L81); [internal/eval/http.go:80–90](../../../../internal/eval/http.go#L80); [internal/eval/model.go:108–112](../../../../internal/eval/model.go#L108).

**Behavior:** Uses pointer usage fields and UsageObserved; records nonnegative values seen in response/events even on failure and sums attempts. Missing usage contributes zero, not inferred tokens.

**CPA comparison:** Original used terminal numeric usage fields and could lose observed usage when stream parsing failed.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native billing still owns actual charges; result token fields and retry aggregation.

**Deletion effect:** Removing preservation loses displayed failed/cancelled usage; it does not disable native gateway accounting. Removing UsageObserved needs callers adjusted.

**Recommendation:** Retain observed usage preservation; distinguish display accounting from native billing.

## C10 — Immutable retries and caller-owned timeout/backoff

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/model.go:47–134](../../../../internal/eval/model.go#L47); [internal/eval/eval.go:289–292](../../../../internal/eval/eval.go#L289); [internal/plugin/codex.go:183–194](../../../../internal/plugin/codex.go#L183); [internal/plugin/model.go:44–85](../../../../internal/plugin/model.go#L44).

**Behavior:** Clones body/headers for each probe attempt, retains exact target, retries status 0/429/500–599 three attempts after 2s/4s, and derives 3-minute probes/10-minute answers from caller context. Slots release before backoff; nil slots allowed for exported ExecuteProbe.

**CPA comparison:** CPA had same attempts/time limits but host body encoding on each call and askCodex created context.Background; >=500 retry lacked 599 upper bound.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C06 error statuses; per-type global semaphores; manager lifetime.

**Deletion effect:** Deleting retries changes existing evaluation robustness/usage; removing child context propagation makes cancellation ineffective. Cloning and nil-slot API support are independent extras.

**Recommendation:** Retain inherited retry counts/timing and timeout limits; retain parent cancellation needed by native runs; review clone/nil-slot convenience separately.

## C11 — Copied Codex client fingerprint and embedded context

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/codex.go:16–55](../../../../internal/eval/codex.go#L16); [internal/eval/codex.go:130–200](../../../../internal/eval/codex.go#L130); [internal/eval/eval.go:276–278](../../../../internal/eval/eval.go#L276); [internal/plugin/codex.go:16–55](../../../../internal/plugin/codex.go#L16).

**Behavior:** Copies tools/base instructions/context, UUID/item construction, User-Agent and protocol headers. Stable installation namespace now derives from sub2api-account:<AccountID>, while each evaluation turn creates fresh thread/turn IDs.

**CPA comparison:** codex.go differs only package/import/comment and removal of CPA askCodex; data/codex_context.json is byte-identical.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** Candy, ModelTrace, Pelican request generation; core/data/codex_context.json in T17 allowlist.

**Deletion effect:** Deleting context/client construction changes preserved prompt/wire fingerprint and breaks turn evaluators. Namespace helper removal needs request callers changed.

**Recommendation:** Retain preserved Codex wire content; approve the account-based namespace explicitly as the native identity adaptation.

## C12 — Candy scoring and result behavior

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/candy.go:8–54](../../../../internal/eval/candy.go#L8); [internal/plugin/candy.go:100–123](../../../../internal/plugin/candy.go#L100).

**Behavior:** Preserves prompt, standalone-21 scoring, 4000-rune saved answer and token/duration fields. Wrong answer remains completed/ok=false; errors now also carry failed/cancelled status and returned error.

**CPA comparison:** Prompt and hasStandalone21 bytes are unchanged apart from shifted lines.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C02; C10; C18 Candy batch manager; original Markdown answer UI.

**Deletion effect:** Deleting scoring/prompt removes Candy functionality; removing returned error/status changes manager outcome handling.

**Recommendation:** Retain inherited Candy semantics and required outcome integration.

## C13 — Fingerprint collector and pure attribution resources

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/fingerprint.go:13–27](../../../../internal/eval/fingerprint.go#L13); [internal/eval/fingerprint.go:61–86](../../../../internal/eval/fingerprint.go#L61); [internal/eval/fingerprint.go:101–172](../../../../internal/eval/fingerprint.go#L101); [internal/eval/fingerprint_algorithm.go:1–424](../../../../internal/eval/fingerprint_algorithm.go#L1); [internal/plugin/fingerprint.go:128–231](../../../../internal/plugin/fingerprint.go#L128).

**Behavior:** Preserves quick/standard/strict sizes, randomized probes, low effort, 2-default/6-max workers, six global slots, stop after eight consecutive request errors, cancellation drain and pure attribution. Adds exported modes/summary, parent cancellation result/error and synchronous progress publication; attributes using effective plain model.

**CPA comparison:** Algorithms differ solely package declaration; all data byte-identical. Threshold and global limits predate native migration.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C04/C10; data/fingerprint_probes.json and baselines; callbacks C16.

**Deletion effect:** Deleting collector removes Fingerprint. Removing new progress/parent status affects manager/UI; algorithms/data need no behavioral simplification.

**Recommendation:** Retain inherited collection/scoring/data; review added public helpers and error-return plumbing against actual UI needs.

## C14 — ModelTrace collector, partial outcomes and golden-bank preservation

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/modeltrace.go:12–44](../../../../internal/eval/modeltrace.go#L12); [internal/eval/modeltrace.go:56–111](../../../../internal/eval/modeltrace.go#L56); [internal/eval/modeltrace.go:113–196](../../../../internal/eval/modeltrace.go#L113); [internal/eval/modeltrace_algorithm.go:1–319](../../../../internal/eval/modeltrace_algorithm.go#L1).

**Behavior:** Preserves three distinct randomized 292–332-number challenges, default effort, 1–3 workers, six global slots, 64000-byte accepted answer bound, sample order and scoring/bank revision. Adds Done/Total/Valid/Errors, parent deadline-aware cancellation counts, partial/failure returned outcome, and exported Summary.

**CPA comparison:** modeltrace_algorithm.go differs solely package; data/modeltrace tree byte-identical.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C04/C10; ModelTrace bank/license; detail UI; C16 callbacks.

**Deletion effect:** Deleting collection/scoring removes ModelTrace. Dropping counters changes progress JSON consumers; deleting bank/license breaks embed.

**Recommendation:** Retain inherited prompts/bank/scoring; review extra progress/return plumbing only as native API additions.

## C15 — Pelican generation, extraction and detail summaries

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/pelican.go:12–81](../../../../internal/eval/pelican.go#L12); [internal/plugin/pelican.go:75–124](../../../../internal/plugin/pelican.go#L75).

**Behavior:** Preserves no-tools prompt, HTML code-block/bare-markup extraction, 256 KiB HTML bound and duration/tokens. Adds failed/cancelled ResultInfo/error returns and Summary method omitting HTML for list responses.

**CPA comparison:** Extraction and prompt unchanged; original handlers already loaded full HTML separately.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C10; native list/detail handler; C31 preview.

**Deletion effect:** Deleting prompt/extraction removes Pelican; removing summary enlarges list payload unless handler revised, not an algorithm change.

**Recommendation:** Retain inherited generation/extraction and bound; approve summary/explicit outcome additions for native list/detail integration.

## C16 — Copied progress callback boundary

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/eval.go:33–47](../../../../internal/eval/eval.go#L33); [internal/eval/eval.go:137–200](../../../../internal/eval/eval.go#L137); [internal/eval/runs.go:428–447](../../../../internal/eval/runs.go#L428).

**Behavior:** Adds JSON-excluded OnProgress callback with copied serial snapshots, initial/final states and post-worker final counters; manager stores progress under stateMu. Candy folds batch counters into a total snapshot.

**CPA comparison:** CPA mutated package-global progress under mu directly.

**Requirement:** **ui** (see provenance key above; current approval remains pending).

**Dependencies:** Collectors C13/C14; manager/UI polling; callbacks must return promptly and not panic/wait.

**Deletion effect:** Deleting callbacks removes live sample progress while final results can still exist; requires manager and UI contract edits.

**Recommendation:** Retain user-visible progress if required; callback shape and constraints are implementation choices pending approval.

## C17 — Managed target admission and account exclusion

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/runs.go:14–55](../../../../internal/eval/runs.go#L14); [internal/eval/runs.go:201–299](../../../../internal/eval/runs.go#L201); [internal/plugin/runs.go:9–57](../../../../internal/plugin/runs.go#L9).

**Behavior:** Resolves display metadata for exact target, rejects nonpositive IDs, deduplicates/excludes by AccountID across groups/types, separates busy and unavailable arrays, and creates failed records for unavailable targets. No client model catalog skips or All wildcard.

**CPA comparison:** CPA operated selectedCredentials and caller model catalog; native resolver is external and cannot substitute.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native TargetResolver; C19 state; C28 UI catalog translation.

**Deletion effect:** Deleting manager admission breaks starts and per-account exclusion. Restoring client model filtering would add a removed behavioral path.

**Recommendation:** Retain exact-target admission/exclusion; explicitly approve unavailable-history versus legacy skipped results.

## C18 — Candy batches, execution settlement and result persistence

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/runs.go:141–162](../../../../internal/eval/runs.go#L141); [internal/eval/runs.go:412–503](../../../../internal/eval/runs.go#L412).

**Behavior:** Runs Candy up to ten times serially per account, different targets concurrently; creates per-iteration run IDs, persists each result and completed count, stops on cancellation/storage failure and releases busy ownership only after work exits.

**CPA comparison:** CPA saved each Candy result but had no durable active/completed intent.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** C17 start; C19 durable intents; C16 progress; StateStore.

**Deletion effect:** Deleting worker disables all evaluations; removing cap/serial ordering changes original Candy semantics. Completed-count persistence is independently new.

**Recommendation:** Retain inherited batch cap/order and settlement; durable per-iteration bookkeeping remains a pending design choice.

## C19 — Persisted active intents and restart recovery

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/state.go:14–38](../../../../internal/eval/state.go#L14); [internal/eval/runs.go:100–138](../../../../internal/eval/runs.go#L100); [internal/eval/runs.go:268–331](../../../../internal/eval/runs.go#L268).

**Behavior:** Adds schema_version=1 and ActiveIntent; saves complete start intent before acknowledging/spawning model workers; restarts append synthetic interrupted records and remove active intents atomically, with protected load/recovery failures.

**CPA comparison:** CPA persisted only completed history and did not recover in-flight runs.

**Requirement:** **history** (see provenance key above; current approval remains pending).

**Dependencies:** C17/C18; StateStore; C20 validators; native history.json location.

**Deletion effect:** Deleting active intent/recovery loses interrupted-run records and pre-ack persistence, but finished-history preservation can remain. Must revise manager/state consumers jointly.

**Recommendation:** Submit active-intent/recovery as its own approval item; historical history requirement does not authorize this full mechanism automatically.

## C20 — Strict history/active schema validation and retention

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/state.go:159–287](../../../../internal/eval/state.go#L159); [internal/plugin/state.go:27–79](../../../../internal/plugin/state.go#L27).

**Behavior:** Validates schema, positive matching account/group, run IDs, model requests, outcome enum, collection counters and active-account uniqueness/options; normalizes empty maps and trims latest 20 per account/type. Invalid data blocks manager start/clear without overwrite.

**CPA comparison:** Original protected unreadable/null JSON and retained20, without versioned/native identity/counter checks.

**Requirement:** **history** (see provenance key above; current approval remains pending).

**Dependencies:** C02/C03/C19; file Load/Save; admission storage latch.

**Deletion effect:** Removing schema/identity validators allows previously blocked data; dropping retention changes inherited history size; changing schema without migration can make current history unreadable.

**Recommendation:** Retain latest-20 preservation; review individual validation rules/protected blocking rather than assuming all are requested.

## C21 — Private atomic file store and no-symlink path checks

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/state.go:46–157](../../../../internal/eval/state.go#L46); [internal/plugin/state.go:86–107](../../../../internal/plugin/state.go#L86).

**Behavior:** Checks every path component for symlink/non-directory/nonregular targets, creates/chmods parent0700, writes a uniquely named0600 temp, rechecks paths, renames atomically. Does not fsync or coordinate processes.

**CPA comparison:** Original used statePath+'.tmp'0600 then rename, without component checks or private directory creation.

**Requirement:** **history** (see provenance key above; current approval remains pending).

**Dependencies:** Native runtime creates store under owned codex-eval directory; C20 schema.

**Deletion effect:** Deleting storage loses restart history; deleting specific path/permission checks relaxes rejection and ownership handling. Direct rename alone preserves original atomic snapshot pattern.

**Recommendation:** Retain atomic persistence; approve extra path checks/parent chmod individually, including their failure effects.

## C22 — Snapshot cloning, record retrieval, clear and flush

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/state.go:289–341](../../../../internal/eval/state.go#L289); [internal/eval/runs.go:505–636](../../../../internal/eval/runs.go#L505); [internal/plugin/state.go:109–159](../../../../internal/plugin/state.go#L109).

**Behavior:** Deep-copies pointer/slice-rich results and snapshots; serves account/run-ID record lookup including Candy; clear supports one/all types and commits only after Save; Flush can clear nonprotected storage error. State shows storage error/quiescing/recovery_blocked intents.

**CPA comparison:** CPA held one mu and returned/cleared global histories; native adds read-safe snapshots and explicit blocked recovery.

**Requirement:** **history** (see provenance key above; current approval remains pending).

**Dependencies:** C20/C21; native admin state/records/clear APIs; runtime deactivation.

**Deletion effect:** Deleting cloning exposes shared mutable state; deleting record/clear removes existing user controls. Flush is used by maintenance, not an extra standalone feature.

**Recommendation:** Retain requested history/detail/clear behavior; review cloning/storage-latch/Flush mechanism as needed for the active consumers.

## C23 — Manager-wide cancellation and shutdown/drain

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/runs.go:638–688](../../../../internal/eval/runs.go#L638); [internal/eval/runs.go:77–96](../../../../internal/eval/runs.go#L77); [internal/plugin/runs.go:60–95](../../../../internal/plugin/runs.go#L60).

**Behavior:** Cancels by run ID across all evaluation kinds, marks cancelling, derives workers from manager lifetime and shutdown closes admission, cancels workers and waits for zero runs plus pending starts. Caller timeout can end wait while workers continue settling.

**CPA comparison:** Original user cancel only Fingerprint/ModelTrace; Candy/Pelican cancellation is newly supported in API although UI retains original buttons.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** C10 context propagation; C18 worker completion; native deactivate/maintenance.

**Deletion effect:** Deleting cancels/drain can leave requests running during cleanup; removing pending-start accounting changes acknowledgment race behavior.

**Recommendation:** Retain cancellation/lifecycle behavior; exact pending-start/drain mechanism requires separate approval.

## C24 — Embedded static UI shell and footer version

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/ui.go:10–36](../../../../internal/eval/ui.go#L10); [internal/eval/web/ui.html:148–165](../../../../internal/eval/web/ui.html#L148); [internal/plugin/plugin.go:79–102](../../../../internal/plugin/plugin.go#L79).

**Behavior:** Caches embedded page/vendor/style/script shell, injects evaluation config/prompts/license and escapes caller nonce. Footer substitutes literal version1 instead of CPA pluginVersion; reports Sub2API title/branding.

**CPA comparison:** Original same shell substitution used actual pluginVersion and no nonce.

**Requirement:** **ui** (see provenance key above; current approval remains pending).

**Dependencies:** Native UI handler CSP nonce; T17 resource mapping; C30/C31.

**Deletion effect:** Deleting embed/shell removes page. Deleting literal footer version only changes display and is independent of functional UI.

**Recommendation:** Retain embedded UI capability; recommend removing misleading literal version display unless user requests it.

## C25 — Native session auth and legacy-shaped API translation

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/app.js:25–70](../../../../internal/eval/web/app.js#L25); [internal/eval/web/ui.html:16–18](../../../../internal/eval/web/ui.html#L16).

**Behavior:** Reads auth_token on every call, rejects absent/401/403 session, no separate management-key input. Translates old run/cancel/DELETE/record paths to native start/cancel/clear/records and native code/data envelope; rejects redirects and keeps30s requests.

**CPA comparison:** CPA decoded management key from obfuscated panel storage and allowed manually entered session key.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native admin API/auth; C26 activation; C27 target state; legacy four-tab handlers.

**Deletion effect:** Deleting auth or translation breaks page operations. Removing legacy adapter layer requires all callers revised; hidden disabled login form is not functional auth.

**Recommendation:** Retain native session reuse/API contract; review compatibility translation as an implementation convenience and hidden login controls separately.

## C26 — First activation controls and readiness gating

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/ui.html:33–34](../../../../internal/eval/web/ui.html#L33); [internal/eval/web/app.js:148–156](../../../../internal/eval/web/app.js#L148); [internal/eval/web/app.js:328](../../../../internal/eval/web/app.js#L328); [internal/eval/web/app.js:440–446](../../../../internal/eval/web/app.js#L440).

**Behavior:** Shows activation state/button, invokes native activate, refreshes targets/status afterward and blocks model-run buttons until ready; repair_required disables activation.

**CPA comparison:** New native initialization UX; original CPA had no separate activate step.

**Requirement:** **ui** (see provenance key above; current approval remains pending).

**Dependencies:** Native status/activate and owned identity; C28 refreshCatalog.

**Deletion effect:** Deleting controls removes requested first-web-activation entry and readiness feedback; native API stays separately callable.

**Recommendation:** Retain explicit first activation flow; individual states/error copy remain implementation details.

## C27 — Group selection and native account/history state adapter

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/app.js:31–32](../../../../internal/eval/web/app.js#L31); [internal/eval/web/app.js:175–189](../../../../internal/eval/web/app.js#L175); [internal/eval/web/app.js:508–513](../../../../internal/eval/web/app.js#L508).

**Behavior:** Persists selected group, renders group accounts with native active/schedulable flags, maps account-keyed histories and active progress to legacy row fields; group change clears selections/catalog and refreshes. History map is account-wide even when selected group differs.

**CPA comparison:** CPA credential rows did not require explicit native group IDs.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Native targets/state endpoints; C02 group metadata; C17 busy account scope.

**Deletion effect:** Deleting target selector removes explicit group selection; deleting state adapter prevents result/progress rows. Group filtering of history is not currently performed.

**Recommendation:** Retain numeric group/account selection and state mapping; explicitly review account-wide history visibility.

## C28 — Native model catalog aggregation and cache

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/catalog.js:3–61](../../../../internal/eval/web/catalog.js#L3); [internal/eval/web/catalog.js:94–137](../../../../internal/eval/web/catalog.js#L94); [internal/eval/web/app.js:78–84](../../../../internal/eval/web/app.js#L78).

**Behavior:** Reads target/status then requests models for every selected-group account, unions IDs and hides image models. Per-account model cache retains inherited identity signature, TTL, dedup/retry handling; refresh failures (except auth) contribute empty lists. No /v1/models or CPA credential sync.

**CPA comparison:** Replaces CPA config/keys inventory and global public models endpoint; exact native aggregation has no explicit separate user request.

**Requirement:** **preserve** (see provenance key above; current approval remains pending).

**Dependencies:** Native account models API; C27; sha256Hex for session-bound cache.

**Deletion effect:** Deleting discovery can leave all model selectors empty. Cache removal changes network volume/prefs but not native model authorization.

**Recommendation:** Retain required native model discovery; approve inherited cache complexity and all-account aggregation separately.

## C29 — Unused pre-start per-account catalog request

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/app.js:249–255](../../../../internal/eval/web/app.js#L249); [internal/eval/web/app.js:39–43](../../../../internal/eval/web/app.js#L39); [internal/eval/web/catalog.js:111–157](../../../../internal/eval/web/catalog.js#L111).

**Behavior:** Before start, verifies global model list then fetches per-account catalogs with six parallel workers and writes body.model_catalog; api immediately deletes model_catalog before sending. Native manager never receives those support facts.

**CPA comparison:** CPA sent this catalog to startRuns for skip behavior; native migration discards it.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** C28 model cache; inherited modelCatalog helpers; API run translation.

**Deletion effect:** Removing the extra modelCatalog call eliminates ignored native-account lookups without changing the transmitted start body; helpers/cache usage should be reviewed separately.

**Recommendation:** Recommend removing the ignored per-account pre-start request as one concrete simplification item; keep global model discovery pending C28.

## C30 — Readonly account switch and obsolete credential machinery

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/components.js:145–154](../../../../internal/eval/web/components.js#L145); [internal/eval/web/app.js:23](../../../../internal/eval/web/app.js#L23); [internal/eval/web/app.js:210](../../../../internal/eval/web/app.js#L210); [internal/eval/web/app.js:341](../../../../internal/eval/web/app.js#L341); [internal/eval/web/credentials.js:60–64](../../../../internal/eval/web/credentials.js#L60).

**Behavior:** Always disables account enable switches with native-account-management hint, leaves switching map/merge/read branches and stub click handler, and retains unused previewCredential after CPA configuredCredentials removal.

**CPA comparison:** CPA actually toggled auth-files/config keys and synthesized hashed credentials; native does neither.

**Requirement:** **ui** (see provenance key above; current approval remains pending).

**Dependencies:** Row display retains status; original checkbox listeners; C27 native accounts.

**Deletion effect:** Removing dead map/preview helper does not change native account mutation because no code fills switching. Removing switch display changes existing UI control appearance.

**Recommendation:** Retain native account status ownership; recommend separately removing unused switch bookkeeping/preview helper. No new account toggle implementation is proposed.

## C31 — Nonce stylesheet conversion for renderer output

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/components.js:3–32](../../../../internal/eval/web/components.js#L3); [internal/eval/web/candy.js:18–21](../../../../internal/eval/web/candy.js#L18); [internal/eval/web/modeltrace.js:54–61](../../../../internal/eval/web/modeltrace.js#L54); [internal/eval/web/app.js:142–145](../../../../internal/eval/web/app.js#L142).

**Behavior:** Converts generated style attributes into deduplicated CSS classes/rules in a nonce style sheet before rows/Markdown/math/rank insertion, avoiding inline-style permission.

**CPA comparison:** CPA permitted unsafe-inline script/style and inserted markup directly.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** Native UI handler nonce CSP; ui.go token replacement; inherited Temml and ModelTrace bars.

**Deletion effect:** Deleting converter while keeping native nonce CSP can lose math layout/ranking style. Deleting only map/dedup changes rule count, not required nonce behavior.

**Recommendation:** Review this added CSP-driven conversion explicitly; retain only if retaining native CSP behavior, with native inventory dependencies.

## C32 — Pelican HTML rewriting under nonce CSP

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/pelican.js:27–62](../../../../internal/eval/web/pelican.js#L27); [internal/plugin/web/pelican.js:27–30](../../../../internal/plugin/web/pelican.js#L27).

**Behavior:** Regex-rewrites opening attributes without altering comments/script raw text, replaces model nonce, moves inline CSS to nonce styles, injects nonce into script/style and serializes DOMParser output into inherited allow-scripts sandbox iframe. Freeze script still stops list animations.

**CPA comparison:** Original srcdoc escaped original HTML directly; sandbox/freeze predated migration.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** C15 stored HTML; native nonce CSP; previews/details; DOMParser.

**Deletion effect:** Removing rewriting preserves original HTML but can block model scripts/styles under present CSP; removing sandbox itself would change inherited isolation.

**Recommendation:** Submit generated-document rewrite as its own approval item; retain required Pelican preview behavior without assuming this implementation is mandated.

## C33 — Native light/dark colors and removal of CPA host styles

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/theme.css:1–36](../../../../internal/eval/web/theme.css#L1); [internal/eval/web/app.js:26–29](../../../../internal/eval/web/app.js#L26); [internal/eval/web/style.css:9–15](../../../../internal/eval/web/style.css#L9).

**Behavior:** Uses native light/dark query/storage/system theme with teal/slate tokens, removes CPAMC/CPAMP backgrounds/blob animation/white theme and removes all:revert descendant reset. Page layout/control stylesheet otherwise unchanged.

**CPA comparison:** style.css differs only one removed reset line; themes replace CPA host blocks rather than redesign tabs.

**Requirement:** **ui** (see provenance key above; current approval remains pending).

**Dependencies:** Native iframe parent theme query in researcher inventory; copied UI components.

**Deletion effect:** Deleting theme file breaks CSS variable appearance; restoring old host themes regresses requested adaptation. Removing query/storage behavior independently affects native theme sync.

**Recommendation:** Retain native theme adaptation; precise token values/design inference remain pending item approval.

## C34 — Stranded CPA subscription-plan controls

**Path:** current embedded runtime source. **Approval:** pending current individual approval.

**Sources:** [internal/eval/web/components.js:101–123](../../../../internal/eval/web/components.js#L101); [internal/eval/web/app.js:387–399](../../../../internal/eval/web/app.js#L387); [internal/eval/web/app.js:179–180](../../../../internal/eval/web/app.js#L179).

**Behavior:** Retains CPA plan taxonomy/rank/filter UI, but credentialPlan requires source=auth_files/provider=codex while native state rows always source=accounts; no plan_type is populated. Filters stay hidden and ranks neutral.

**CPA comparison:** Copied CPA feature now unreachable after target adapter change.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** Inherited plan badges/filters and preference controls; native accounts mapping.

**Deletion effect:** Removing unreachable taxonomy/filter helpers needs hidden DOM/render calls adjusted, with no current native plan output lost.

**Recommendation:** Recommend reviewing removal of unreachable CPA plan controls as an individual simplification; do not add a native plan feature.

## T01 — Prebuilt-only installation payload and release constants

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:26–35](../../../../install.py#L26); [install.py:794–807](../../../../install.py#L794); [install.py:896–914](../../../../install.py#L896).

**Behavior:** Selects versioned v0.2.16 amd64 candidate and exact eight-file payload; default installation uses local complete bundle or downloads fixed release archive. No Git/source/Go/frontend build function exists in current installer.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T02 acquisition; T08 compatibility; T16 packaging; TOOL_FACTS used by maintenance/idempotence.

**Deletion effect:** Deleting binary-loading path disables installation. Deleting historical resolver/build files cannot simplify this live installer because it does not call them.

**Recommendation:** Retain user-requested prepared binary replacement; do not describe historical on-host source build as current runtime.

## T02 — Local/HTTPS release lookup and bounded allowlisted extraction

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:148–225](../../../../install.py#L148); [install.py:901–914](../../../../install.py#L901).

**Behavior:** Recognizes only a complete local bundle; otherwise downloads pinned plugin release-tag archive with environment proxies disabled and HTTPS-only redirects. Bounds compressed/uncompressed bytes; checks unique exact member names/types/directories; writes exclusive private extracted files.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T01 names/allowlist; T16 archive format; online GitHub asset path.

**Deletion effect:** Removing online path loses one-command bootstrap; removing exact/size/type checks admits other files. Proxy disabling and size thresholds are implementation-added choices.

**Recommendation:** Retain required local/online acquisition and minimal payload; review exact extraction/proxy bounds individually without inventing new binary attestation.

## T03 — Explicit Docker target choice and immutable selection

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:347–390](../../../../install.py#L347); [install.py:1082–1093](../../../../install.py#L1082).

**Behavior:** Lists all Docker containers via quiet full IDs/JSON inspect, then requires explicit --container or numbered terminal choice even with one candidate. /dev/tty read/write handles piped bootstrap; selection uses immutable ID and rechecks name/ref before effects.

**Requirement:** **selection** (see provenance key above; current approval remains pending).

**Dependencies:** Docker CLI/Python3.9; T04 inspect; all restart/replace effects.

**Deletion effect:** Deleting selection/rechecks risks wrong container mutation; removing quiet-ID compatibility reintroduces reported --format failures.

**Recommendation:** Retain explicit user choice and compatible enumeration; separately review recheck frequency if simplifying.

## T04 — Native container layout/identity preflight

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:399–447](../../../../install.py#L399); [install.py:647–675](../../../../install.py#L647).

**Behavior:** Requires running Linux/container writable rootfs, normal billing mode, writable persistent data mount, safe unmapped executable; derives uid/gid from PID1, resolves executable from proc/startup hints, reads optional -version/image labels and checks safe path ancestors. Initial inspect permits amd64/arm64; fixed candidate later requires amd64.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T05/T07 functional probe; T08 fixed layout/version check; journal/backups rely on runtime uid/gid.

**Deletion effect:** Removing path/ownership discovery can make backup/replacement target or permissions wrong. Removing version hints does not alter native dispatch but affects known-version conflict detection.

**Recommendation:** Retain facts necessary for requested replacement; each additional fail-closed check and PID1 identity assumption remains a concrete pending review.

## T05 — Owned clientless probe image acquisition/cache

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:60–66](../../../../install.py#L60); [install.py:276–342](../../../../install.py#L276).

**Behavior:** Uses alpine/git2.49.1 solely for nc; probes cached inspected Linux image in named no-network temporary helper, cleans exact helper, pulls/retries transient failures once if absent/unusable and pins subsequent use to inspected image ID. Successful selected image remains available for later probes.

**Requirement:** **selection** (see provenance key above; current approval remains pending).

**Dependencies:** T06 nc transport; Docker pull/cache; cleanup retains probe image.

**Deletion effect:** Deleting helper acquisition breaks images lacking curl/wget/nc; deleting transient classifier changes only pull retry behavior, not Sub2API source selection.

**Recommendation:** Retain clientless probing due recorded missing nc failure; review use of this image/cache probe/retry mechanism rather than assuming all are user mandated.

## T06 — Custom nc HTTP framing and helper teardown

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:69–127](../../../../install.py#L69); [install.py:449–546](../../../../install.py#L449).

**Behavior:** Runs exact named nc helper in selected container network namespace and sends only /health or /setup/status. Keeps stdin open, parses bounded HTTP status/encoding/content-length/chunk framing through HTTPResponse plus manual chunk terminator checks, uses selectors/deadline and strict process/helper cleanup preserving raw output.

**Requirement:** **selection** (see provenance key above; current approval remains pending).

**Dependencies:** T05 probe image; T07 endpoint discovery; T10 startup verification.

**Deletion effect:** Deleting framing helper prevents current clientless readiness; removing manual trailer/size handling alters incomplete/oversize acceptance. Teardown changes can leave owned helper resources.

**Recommendation:** Retain working clientless probes; submit parser/bounds/teardown complexity as concrete approval choices, with no replacement transport proposal.

## T07 — Owned listener discovery, hint fallback and functional readiness

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:547–645](../../../../install.py#L547).

**Behavior:** Walks PID1 children at most3 depths/32 processes, matches executable socket inodes to proc TCP/TCP6 listeners, decodes addresses and prioritizes configured port. Private network fallback uses env/exposed port hints; host/shared namespaces require owned sockets. Up to8 endpoints/48s must return exact native health and completed setup fields.

**Requirement:** **selection** (see provenance key above; current approval remains pending).

**Dependencies:** T04 executable; T06 HTTP; health_host/port stored in facts.

**Deletion effect:** Deleting owned discovery can probe wrong shared-network listener or miss nondefault binding; removing exact setup checks broadens accepted applications. Limits/hint fallbacks are assistant-selected.

**Recommendation:** Retain endpoint discovery needed by reported port/client failure; explicitly review bounded proc scanning and exact response conditions.

## T08 — Fixed amd64, /app/data and known-version conflict checks

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:794–807](../../../../install.py#L794); [install.py:916–920](../../../../install.py#L916).

**Behavior:** Requires amd64 and exactly/app/data; rejects only known semver differing0.2.16 while allowing unknown version. Candidate checked for nonempty ordinary payload file; no installed checksum/ELF/build-metadata/image allowlist is reinstated.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T04 hints; T01 release-specific candidate; native build data_dir fact.

**Deletion effect:** Removing architecture/layout restrictions can run wrong artifact/storage path; removing known-version check permits silent downgrade. These checks are distinct from removed binary fingerprint gates.

**Recommendation:** Retain requested release/layout compatibility; present version-conflict policy separately from binary identity checks.

## T09 — Protected journal, private root and original backup

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:392–397](../../../../install.py#L392); [install.py:677–716](../../../../install.py#L677); [install.py:936–953](../../../../install.py#L936).

**Behavior:** Creates owned0700 codex-eval directory, schema3 installation journal0600 and original-sub2api backup0500; checks no-symlink ancestors, exact runtime ownership/modes and supported phase/deployment/path fields. Writes journal/facts via unique stage, chown/chmod/sync/rename.

**Requirement:** **transaction** (see provenance key above; current approval remains pending).

**Dependencies:** T04 uid/gid and paths; T10 phases; native installation ownership templates.

**Deletion effect:** Deleting backup loses reversible replacement; deleting journal or schema/permission checks affects crash/reapply and native identity ownership. Changing state requires explicit treatment of existing files.

**Recommendation:** Retain original backup and needed ownership record; review strict modes/schema/phase validation as additional complexity.

## T10 — Atomic replacement, restart checks and failure restoration

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:809–812](../../../../install.py#L809); [install.py:836–894](../../../../install.py#L836); [install.py:954–993](../../../../install.py#L954).

**Behavior:** Stages candidate beside executable, writes installer facts, renames/restarts only selected container, polls Docker health/native health/maintenance status. On post-replace exception restores original, starts stopped container via protected copied backup or renames live backup, verifies health and keeps failure journal; cleans exact stage.

**Requirement:** **transaction** (see provenance key above; current approval remains pending).

**Dependencies:** T09 backup/journal; T06 probes; T12 safe maintenance; native runtime status.

**Deletion effect:** Deleting compensation leaves failed patched binary; dropping stage/restart checks changes atomicity/readiness; replacement flag deliberately set before potentially failing mv to ensure restoration attempt.

**Recommendation:** Retain required replacement/recovery behavior; each phase/check remains pending individual implementation approval.

## T11 — Patch digest facts and reinstall short circuit

**Path:** packaged installer and build metadata. **Approval:** pending current individual approval.

**Sources:** [install.py:26–29](../../../../install.py#L26); [install.py:921–933](../../../../install.py#L921); [install.py:947–965](../../../../install.py#L947); [adapter/build.py:96–104](../../../../adapter/build.py#L96).

**Behavior:** Stores source revision/plugin source digest in installer facts/journal; matching installed journal digest plus safe maintenance/backup/running checks skips replacement. Build injects source/digest into runtime and writes development build-facts.json. Never compares installed binary bytes.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T17 functional_hash defines digest; native runtime facts consumer; T09 private records.

**Deletion effect:** Removing digest-based short circuit may reinstall each time; dropping specific source/digest facts requires runtime ownership/status contracts adjusted. Deleting build-facts output only affects development provenance consumers.

**Recommendation:** Review metadata/idempotence necessity; do not conflate plugin digest matching with removed host binary attestation.

## T12 — Safe maintenance status, rollback and uninstall

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:813–834](../../../../install.py#L813); [install.py:995–1079](../../../../install.py#L995); [install.py:1082–1099](../../../../install.py#L1082).

**Behavior:** Checks capability flag before maintenance subcommand; requires bounded machine JSON with exact deployment ID/schema/known state. Deactivate must confirm drained and owned-key revocation before original restoration; phase-resumable rollback/uninstall retain history/backup/journal. status/inspect print actual installed versus active state.

**Requirement:** **transaction** (see provenance key above; current approval remains pending).

**Dependencies:** Native maintenance CLI/socket/runtime ownership and Key revoke; T09/T10.

**Deletion effect:** Deleting capability/schema checks can invoke unsupported original binaries or restore before owned cleanup. Removing rollback/uninstall loses requested reversible operation.

**Recommendation:** Retain requested status/rollback/uninstall and exact owned cleanup; review inspect alias/response/phase machinery as implementation details.

## T13 — Raw diagnostics and live subprocess supervision

**Path:** packaged installer. **Approval:** pending current individual approval.

**Sources:** [install.py:130–145](../../../../install.py#L130); [install.py:245–274](../../../../install.py#L245); [install.py:720–792](../../../../install.py#L720); [install.py:1100–1102](../../../../install.py#L1100).

**Behavior:** Relays complete failed Docker stdout/stderr, retains exception chains/traceback; live pulls stream both pipes with30s liveness notice,120s timeout and process-group termination/drain. live_process accepts logfile but the current call graph supplies none.

**Requirement:** **selection** (see provenance key above; current approval remains pending).

**Dependencies:** T05 image pulls; T06 own transport diagnostics; no current builder invocation.

**Deletion effect:** Deleting raw diagnostics violates recorded failure-reporting request; deleting unused logfile branches changes no current installer call; supervision teardown remains needed for pulls.

**Recommendation:** Retain raw error reporting; recommend removing the unused logfile option/branches as a separate narrow simplification.

## T14 — Independent cleanup of attributable stopped helpers/workspaces

**Path:** packaged cleanup. **Approval:** pending current individual approval.

**Sources:** [cleanup.py:14–18](../../../../cleanup.py#L14); [cleanup.py:46–106](../../../../cleanup.py#L46); [cleanup.py:109–130](../../../../cleanup.py#L109); [cleanup.py:138–173](../../../../cleanup.py#L138); [cleanup.py:176–208](../../../../cleanup.py#L176).

**Behavior:** Inventories exact owned helper/buildx/temp prefixes; only inactive unmounted/stopped helpers and attributable owner-held nonsymlink standard-temp workspaces eligible. Optional Buildx detected via root help and conservative plain-table header/driver parsing; active/unknown retained. Exact selected IDs/dry-run, refreshed inventory per removal and symlink-safe rmtree.

**Requirement:** **cleanup** (see provenance key above; current approval remains pending).

**Dependencies:** Docker/proc process inspection; old naming conventions; T05 helper retention; T15 image candidates.

**Deletion effect:** Deleting cleanup violates separately requested utility; deleting ownership/current-use checks broadens removals; no resource is auto-selected without explicit input.

**Recommendation:** Retain requested selective cleanup; exact process/table/ancestry mechanisms remain independently reviewable.

## T15 — Unowned Go/Node image and anonymous tool candidate cleanup

**Path:** packaged cleanup. **Approval:** pending current individual approval.

**Sources:** [cleanup.py:121–137](../../../../cleanup.py#L121); [cleanup.py:151–164](../../../../cleanup.py#L151); [cleanup.py:187–203](../../../../cleanup.py#L187).

**Behavior:** Offers unused unmounted tool-container candidates based on inspected Go/Node tags (including digest-started images) and unreferenced Go/Node images; requires exact selection and refreshed reference/eligibility checks. Other tags/shared/reference ambiguity retained; Alpine probe image is not a tool candidate.

**Requirement:** **cleanup** (see provenance key above; current approval remains pending).

**Dependencies:** T14 inventory/process checks; image JSON inspection; user exact selection.

**Deletion effect:** Deleting generic candidates keeps those older unlabeled build leftovers; deleting reference checks could remove shared resources. This utility never deletes volumes or global prunes.

**Recommendation:** Retain only the explicitly requested old-build candidate scope and refreshed checks, pending item approval.

## T16 — Eight-file packaging, fixed outputs and atomic deterministic archive

**Path:** developer packager current ci. **Approval:** pending current individual approval.

**Sources:** [scripts/package-sub2api.py:15–35](../../../../scripts/package-sub2api.py#L15); [scripts/package-sub2api.py:45–126](../../../../scripts/package-sub2api.py#L45); [scripts/package-sub2api.py:129–143](../../../../scripts/package-sub2api.py#L129).

**Behavior:** Packages install/cleanup/binary/five license texts, checks source/binary parity; fixed canonical raw/archive defaults, automatic system-temp stage or explicit overrides. Writes normalized USTAR/gzip(no timestamp/name/PAX) from exact paths, validates raw member bytes/types and atomically replaces archive.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T17 license/source map and facts helper; T18 binary build; T25 workflow fixed paths.

**Deletion effect:** Deleting package path loses requested artifact; removing parity/tar metadata validation may reintroduce unexpected source/AppleDouble payload. Optional stage/override controls are developer conveniences.

**Recommendation:** Retain exact minimal payload/fixed outputs; review overrides and normalization/parity checks separately if reducing developer tooling.

## T17 — Functional build source allowlist and package-hash coupling

**Path:** developer build and package current. **Approval:** pending current individual approval.

**Sources:** [scripts/universal_payload.py:1–12](../../../../scripts/universal_payload.py#L1); [scripts/prebuilt_dev.py:44–51](../../../../scripts/prebuilt_dev.py#L44); [scripts/package-sub2api.py:12–18](../../../../scripts/package-sub2api.py#L12); [scripts/sub2api_build.py:648–686](../../../../scripts/sub2api_build.py#L648).

**Behavior:** Explicit SOURCES selects core/adapter/runtime resources for developer builds; package only uses its LICENSE entries. functional_hash includes adapter Dockerfile and all functional resources except installer/non-ModelTrace licenses, omits cleanup.py and tests. Package parses bounded single literal TOOL_FACTS via old sub2api_build and requires hash match.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** All C items/native inventory; T11 source digest; T16 packaging; historical helper imports T26.

**Deletion effect:** Deleting entire universal_payload or sub2api_build breaks current package/build even though most older machinery is unused. Changing source hashes requires updating facts together; Dockerfile is hash-active despite not build-executed.

**Recommendation:** Retain actual source allowlist; review source-hash/facts coupling as an implementation addition. Do not remove whole historical modules before addressing live consumers.

## T18 — Developer v0.2.16 source acquisition and build profile

**Path:** developer build current ci. **Approval:** pending current individual approval.

**Sources:** [scripts/prebuilt_dev.py:14–60](../../../../scripts/prebuilt_dev.py#L14); [.github/workflows/release.yml:24–32](../../../../.github/workflows/release.yml#L24).

**Behavior:** Verifies peeled annotated tag equals fixed revision, downloads codeload archive, checks safe archive paths/types and records observed archive SHA; copies SOURCES and builds profile amd64,/app/data,version0.2.16. Profile describes local-go1.27.2/local-node24 but Python does not probe those executables; CI sets Go1.27.2/Node24.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T17 source map/hash; T19 actual builder; upstream revision/network; repository developer README.

**Deletion effect:** Deleting loses currently documented release build; tag/archive/profile checks are developer provenance, not installer compatibility. Removing profile tool labels changes facts only; changing Node requires existing upstream contract review.

**Recommendation:** Retain actual user-requested build from official0.2.16; profile metadata/tool selection extras remain pending. Pristine baseline used upstream Node20 and must not be misrepresented as this builder.

## T19 — Actual frontend/backend builder and generated data-path facts

**Path:** developer build current ci. **Approval:** pending current individual approval.

**Sources:** [adapter/build.py:14–51](../../../../adapter/build.py#L14); [adapter/build.py:53–107](../../../../adapter/build.py#L53).

**Behavior:** Verifies selected source revision/archive provenance, hashes locks, runs native AST rules/runtime copies/UI edits, copies staged core/ModelTrace license, emits EvalDataDir/control.sock init. Frozen pnpm build is embedded, compiler type check runs, static embed Linux target built with release/source/digest ldflags, lock hashes compared and build facts written.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T18 source/profile; T20 AST engine; T21 type-export extra; native structural inventory; T17 allowlist.

**Deletion effect:** Deleting builder disables plugin-enabled compile; removing generated facts breaks/app/data runtime location; source/lock validation and build-facts output are additional decisions beyond compile.

**Recommendation:** Retain required actual frontend/core/native runtime compilation; review source/lock guard and generated build metadata individually.

## T20 — AST structural edit engine and source transaction

**Path:** developer build current. **Approval:** pending current individual approval.

**Sources:** [adapter/main.go:45–95](../../../../adapter/main.go#L45); [adapter/main.go:173–258](../../../../adapter/main.go#L173); [adapter/main.go:260–349](../../../../adapter/main.go#L260); [adapter/build.py:17–26](../../../../adapter/build.py#L17).

**Behavior:** Matches exact declaration/scoped canonical AST sequences with neighboring anchors, fails missing/ambiguous/overlapping hooks, stages formatted Go then writes all required edit files. Operates on current rules.json, not historical *.patch files; engine compiled from main.go/typecheck.go only.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** Researcher inventory owns exact individual native rules; T19 invokes engine; added imports/declarations in recipe.

**Deletion effect:** Deleting engine prevents applying current native source hooks. Removing strict matching changes source acceptance; it does not remove a specific runtime behavior by itself.

**Recommendation:** Retain only because current build requires this application mechanism, pending approval; judge actual runtime rules separately and do not retain universal compatibility guards merely by prior spec.

## T21 — Separate native compiler-export contract gate

**Path:** developer build current. **Approval:** pending current individual approval.

**Sources:** [adapter/typecheck.go:15–95](../../../../adapter/typecheck.go#L15); [adapter/typecheck.go:99–189](../../../../adapter/typecheck.go#L99); [adapter/build.py:92–93](../../../../adapter/build.py#L92).

**Behavior:** Before final build, go list -deps -export compiles graph, loads go/types exports and requires a large named method/field map including receiver identity and internal native types. This duplicates some final compiler compatibility checks while asserting extra expected type shapes.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T20 types operation; native service recipe/runtime; requires compiled dependency graph.

**Deletion effect:** Removing this invocation/map leaves final Go compilation but loses explicit named API/field-shape gate and its failures; historical type-contract tests may need consumer adjustment.

**Recommendation:** Recommend separate review of removing this additional gate; no specific user request for compiler-export map was found. Do not add a substitute acceptance gate.

## T22 — Development-only AST recipe exporter

**Path:** developer unused current build. **Approval:** pending current individual approval.

**Sources:** [adapter/export_dev.go:136–194](../../../../adapter/export_dev.go#L136); [adapter/build.py:23–25](../../../../adapter/build.py#L23).

**Behavior:** Optional extract operation derives rules by comparing before/after Go ASTs and serializes recipe. build.py compiles only main.go/typecheck.go, so export_dev.go is absent from current build engine and SOURCES.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** Manual developer command may compile exporter with engine; current rules remain ordinary files.

**Deletion effect:** Deleting exporter does not change current binary/package/install path; removes maintenance regeneration tool only.

**Recommendation:** Recommend removing unused exporter only through its own approval; no runtime justification exists.

## T23 — Test materialization overlays and source parity

**Path:** developer test current ci. **Approval:** pending current individual approval.

**Sources:** [scripts/prebuilt_test.py:13–36](../../../../scripts/prebuilt_test.py#L13); [scripts/universal_dev.py:19–39](../../../../scripts/universal_dev.py#L19); [tests/go.mod:1–4](../../../../tests/go.mod#L1); [scripts/export-sub2api-overlays.py:7–20](../../../../scripts/export-sub2api-overlays.py#L7).

**Behavior:** Re-extracts verified original build archive into temp Git baseline, reapplies current structural adapters, copies v0.2.15-named test-overlay patches plus tests/overlays/codexeval, and checks existing backend/frontend bytes unchanged. Separate tests/go.mod hides raw host-dependent overlays from root Go discovery. Overlay exporter derives add-only test patches from tests/host-overlay.

**Requirement:** **native** (see provenance key above; current approval remains pending).

**Dependencies:** Current native/core regressions; test-overlay path is active despite versioned name; T20; T24 gates.

**Deletion effect:** Deleting test overlays/materialization removes actual focused native/core regressions; deleting versioned production patches does not remove these tests. Exporter is dev-only but maintains active overlay files.

**Recommendation:** Retain necessary existing regression sources per user instruction; review materialization/export helpers individually, without treating old path name as an unused patch indication.

## T24 — Current local gates, historical helper tests and legacy UI test owner

**Path:** developer tests current ci. **Approval:** pending current individual approval.

**Sources:** [scripts/test-sub2api.sh:20–41](../../../../scripts/test-sub2api.sh#L20); [scripts/universal_dev.py:41–47](../../../../scripts/universal_dev.py#L41); [tests/installer/current_suite.py:1–26](../../../../tests/installer/current_suite.py#L1); [scripts/test_sub2api_build.py:14–21](../../../../scripts/test_sub2api_build.py#L14); [scripts/test-ui.cjs:7](../../../../scripts/test-ui.cjs#L7).

**Behavior:** Current check runs historical build-helper test module, selectively current installer suite, facts check, shell syntax, root vet/test and scripts/test-ui.cjs. That UI harness reads internal/plugin/web (CPA), not internal/eval/web. Native gates run vet/unit/integration in adapted workspace; race invocation already removed by main.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T17 hash helpers; T23 overlays; CPA root tests remain retained; tests/ui/sub2api-browser.js is separate actual native browser fixture.

**Deletion effect:** Deleting whole check would remove required regression coverage; dropping historical helper tests changes only historical functionality evidence. Passing CPA UI harness cannot prove current native UI.

**Recommendation:** Review additional historical helper/shell/UI-gate inclusion separately; retain applicable upstream/current user-approved checks. No new gates or tests are proposed.

## T25 — Sub2API CI/release delivery wiring

**Path:** developer ci release. **Approval:** pending current individual approval.

**Sources:** [.github/workflows/release.yml:12–61](../../../../.github/workflows/release.yml#L12); [.github/workflows/release.yml:63–77](../../../../.github/workflows/release.yml#L63).

**Behavior:** On tags/workflow dispatch verifies plugin tag, builds patched frontend/binary, runs frontend tests/local/native gates, packages fixed archive, runs owned disposable official-image lifecycle, uploads archive; tag job publishes only archive after delivery. Publication has not been run in this audit.

**Requirement:** **prebuilt** (see provenance key above; current approval remains pending).

**Dependencies:** T18/T16/T23/T24; Docker lifecycle fixture; TOOL_FACTS release tag.

**Deletion effect:** Deleting workflow removes automated build/release route; deleting individual extra steps changes release requirements. No runtime/installer behavior changes until a later artifact build.

**Recommendation:** Retain requested release artifact wiring pending review of each extra check; do not infer current permission to publish from workflow contents.

## T26 — Historical pinned-host manifest and unused source/patch machinery

**Path:** historical developer code partially live import. **Approval:** pending current individual approval.

**Sources:** [compatibility/sub2api-v0.2.15.json:1–27](../../../../compatibility/sub2api-v0.2.15.json#L1); [compatibility/sub2api-v0.2.15.json:65–114](../../../../compatibility/sub2api-v0.2.15.json#L65); [scripts/sub2api_build.py:77–152](../../../../scripts/sub2api_build.py#L77); [scripts/sub2api_build.py:163–224](../../../../scripts/sub2api_build.py#L163); [scripts/sub2api_build.py:266–391](../../../../scripts/sub2api_build.py#L266); [scripts/sub2api_build.py:445–645](../../../../scripts/sub2api_build.py#L445); [scripts/sub2api_build.py:689–845](../../../../scripts/sub2api_build.py#L689).

**Behavior:** Old manifest pins0.2.15 images/platforms/toolchain/core/production patches; old helpers verify source/tools, parse/validate patch paths, prepare/reuse/backup workspaces, apply production patches, stage core, materialize old tests and write provenance/CLI commands. Current prebuilt build/install do not call these functions. Active imports still use _sha256_file, parse_release_facts and validate_installer_facts.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T17/T16 live small helpers; scripts/pnpm and historical foundation tests still consume old manifest/CLI.

**Deletion effect:** Deleting whole file immediately breaks current packaging/CI; deleting only unused helper families leaves current runtime unchanged but affects historical tests/pnpm CLI.

**Recommendation:** Recommend reviewing unused historical helper families individually after identifying live imports; no bulk module deletion.

## T27 — Historical versioned production and combined patch trees

**Path:** historical not current runtime input. **Approval:** pending current individual approval.

**Sources:** [patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch:1–14](../../../../patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch#L1); [patches/sub2api/v0.2.15/005-codex-eval-deployment-evidence.patch:1–2](../../../../patches/sub2api/v0.2.15/005-codex-eval-deployment-evidence.patch#L1); [compatibility/sub2api-v0.2.15.json:65–67](../../../../compatibility/sub2api-v0.2.15.json#L65); [adapter/build.py:17–38](../../../../adapter/build.py#L17); [scripts/universal_dev.py:30–31](../../../../scripts/universal_dev.py#L30).

**Behavior:** Top-level001–006 combined patches and production001–004/007/008 represent earlier0.2.15 wiring/topology/loopback evidence. Current builder reads structural rules/current runtime/core, not these production/combined patches. Versioned test-overlay directory is separately live under T23.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** Old manifest/helper tests T26 reference production tree; pristine history retained separately.

**Deletion effect:** Deleting historical production/combined files does not alter current build/runtime/install, but can break old manifest helper tests/evidence links. Deleting test-overlay instead would affect active regressions.

**Recommendation:** Recommend review of each obsolete historical patch family for removal; preserve evidence and current test-overlay distinction.

## T28 — Retired wrappers and stale tool-shim contracts

**Path:** historical developer entries. **Approval:** pending current individual approval.

**Sources:** [scripts/build-sub2api.sh:1–10](../../../../scripts/build-sub2api.sh#L1); [scripts/export-sub2api-production.py:1–7](../../../../scripts/export-sub2api-production.py#L1); [scripts/generate-wire.sh:4–26](../../../../scripts/generate-wire.sh#L4); [scripts/pnpm:4–7](../../../../scripts/pnpm#L4); [scripts/pnpm-support.sh:1–6](../../../../scripts/pnpm-support.sh#L1); [scripts/test-sub2api.sh:25–30](../../../../scripts/test-sub2api.sh#L25).

**Behavior:** build-sub2api.sh check still validates facts; normal invocation prints obsolete install-time-build instruction then exit2. Production exporter is explicit exit2 tombstone; Wire temp-mod helper and pnpm compatibility0.2.15 shim/support have no current adapter build caller, although shell syntax gate still reads them.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T26 historical manifest/CLI; test-sub2api shell syntax; current build uses corepack directly and existing wire_gen recipe.

**Deletion effect:** Deleting wrappers/shims leaves current Python build/install/runtime unchanged except shell syntax/reference adjustments; deleting check alias loses a convenience facts command.

**Recommendation:** Recommend separate removal review for each retired/stale entry, especially misleading automatic-build output; do not invoke them.

## T29 — Unused Docker build recipe still enters functional digest

**Path:** developer unexecuted hash input. **Approval:** pending current individual approval.

**Sources:** [adapter/Dockerfile:1–14](../../../../adapter/Dockerfile#L1); [scripts/universal_payload.py:5–12](../../../../scripts/universal_payload.py#L5); [scripts/prebuilt_dev.py:49–53](../../../../scripts/prebuilt_dev.py#L49).

**Behavior:** Docker multi-stage recipe copies source/core/adapter/licenses and calls build.py with Go/Node image ARGs; current prebuilt_dev directly runs local Python builder and installer has no docker build. SOURCES still copies Dockerfile and functional_hash includes it.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T17 digest/facts; old on-host build history; no current CI Docker build.

**Deletion effect:** Deleting Dockerfile alone changes functional hash and breaks facts validation even though executable runtime source is unchanged; must address map/facts together after approval.

**Recommendation:** Recommend removal review of this unexecuted hash input, preserving awareness of live digest dependency.

## T30 — Obsolete two-release universal CI runner

**Path:** historical developer unused. **Approval:** pending current individual approval.

**Sources:** [tests/installer/universal_ci.py:17–44](../../../../tests/installer/universal_ci.py#L17).

**Behavior:** Packages current minimal prebuilt release then tries universal_dev.prepare using that release as a source package for two latest upstream tags; source package no longer contains adapter/core, so path is stale. Not referenced by current workflow/current_suite.

**Requirement:** **none** (see provenance key above; current approval remains pending).

**Dependencies:** T23 adapt needs source package; T16 now only emits eight files; historical universal deliverable superseded.

**Deletion effect:** Deleting runner does not change current CI/runtime/install; loses historical entry only. No current execution/reproduction performed.

**Recommendation:** Recommend reviewing removal of stale runner; do not interpret it as current universal support or run it during planning.

## T31 — Historical installer source-progress/build tests with live selective imports

**Path:** historical tests mixed with current suite. **Approval:** pending current individual approval.

**Sources:** [tests/installer/current_suite.py:3–24](../../../../tests/installer/current_suite.py#L3); [tests/installer/test_preflight.py:1–24](../../../../tests/installer/test_preflight.py#L1); [scripts/test_sub2api_build.py:14–21](../../../../scripts/test_sub2api_build.py#L14).

**Behavior:** Old source-builder/progress/universal test files remain maintained. Current_suite selects transaction/clientless/prebuilt/cleanup/compat cases plus named preflight/discovery cases from old modules; foundation test module remains invoked by current local gate.

**Requirement:** **cleanup** (see provenance key above; current approval remains pending).

**Dependencies:** T24 selected tests; historical source/progress test modules; user earlier retained maintained tests.

**Deletion effect:** Deleting all old-looking test files can break live imports/selected tests and discard deliberately retained sources. Dropping invocation of obsolete-only cases is different from deleting tests.

**Recommendation:** Retain test source unless explicitly approved individually; review gate membership separately. No test deletion or assertion weakening is authorized.

## Handoff

Start with one concrete item. C29 is a narrow ignored pre-start model-catalog request; T21 is an extra compiler-export gate; T28 covers retired/stale developer entries with live syntax references; T29 is an unexecuted Dockerfile with live digest coupling. These are proposals, not a batch removal order. Preserve test sources and diagnostic evidence until specifically decided.

The native inventory must resolve cross-owner dependencies before any approved core/installer edit is implemented. Runtime templates, exact native rule edits, identity/auth/scheduler/menu/maintenance behavior are intentionally not duplicated here.


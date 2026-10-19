# Native Browser Check

## Candidate and Fixture

The candidate source is `/tmp/upstream-sync-native/source`.
The temporary fixture is `/tmp/sub2api-upstream-sync-preview.cjs`.
It serves the native Vue layout and candidate evaluation iframe on port 60232.
Chrome used synthetic authentication, groups, accounts, settings, and result records.
No production account or supplier request was used.

## Results

Native sidebar navigation opened the evaluation iframe with the saved Settings tab.
The delayed settings response completed before shared group loading.
Switching directly to fingerprint selected group B and showed B-only and shared accounts.
No manual refresh was required.

The page showed Candy, fingerprint, Pelican, and Settings.
The fingerprint toolbar retained the native model dropdown and per-account concurrency, with default 3.
No retired mode control or ModelTrace tab remained.
The actual dark-theme render retained the native shell, account table, and controls.
No explanatory small print or principle panel was added.
The per-account concurrency control accepted a change from 3 to 2.

A partial record showed the model comparison and a dash for null probability.
Its detail showed candidate scores and dashes for null probabilities.
The challenge section showed attempt counts, acceptance, prompts, instructions, and answers.
Expanding a challenge worked.

A `not_model` record showed `不像语言模型回答` in the row and detail.
Its detail omitted candidate rankings and probabilities.

The first synthetic ranking omitted its required score field.
The fixture was corrected and restarted. The final detail showed numeric scores.
This fixture correction did not change product code.

The Pelican page retained the selected group and both account rows.
Its fixed one-request caption conflicted with the new retry behavior.
The implementation removed that caption without replacement text.
After the fixture restarted, the final browser tree confirmed its absence and loaded both account rows.
The final screenshot request timed out. Earlier native-page and detail screenshots supplied the render evidence.

## Scope

The existing UI harness owns requests, cancellation IDs, custom-model behavior, and retired-tab fallback.
This check inspects the actual changed render and Settings-to-test navigation.
The README uses the visible tab label `指纹`.

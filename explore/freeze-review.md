# Freeze review candidate — 2026-10-29

This is an unfrozen candidate, not an operator approval or a holdout result.
The authoritative machine-readable candidate is `records/moneyzg-five-minute/spec.json`
(SHA-256 `14d326cf92bb92fefd5ef65fcaa6fe8e074aeea75aca3f9073df86ad1b44697c`),
which binds the source manifest. The study runner rejects a lineage with a
different specification hash. The manifest includes verified development
identities; the reserved holdout identities remain pending acquisition.

## Contract and boundaries

All intervals are UTC and half-open. Development evaluation is
[2024-01-01 00:00, 2024-12-30 00:00), preceded by 500 bars beginning
2023-12-30 06:20. Holdout warm-up is [2024-12-30 06:20, 2025-01-01 00:00),
exactly 500 five-minute bars. Evaluation is [2025-01-01 00:00, 2026-01-01
00:00): 105,120 bars and 365 days. No pivot outside evaluation is eligible;
warm-up bars seed indicators but do not generate eligible trades. No open
position is forcibly liquidated: a position still open at the boundary makes
the run INCOMPLETE, with no scientific verdict. Completed trades are assigned
to their UTC exit day; all days, including zero-trade days, are retained.
This reserved historical holdout is prospective relative to the freeze, not
calendar-forward validation.

The complete mechanical and funding rules are in the specification. In brief,
close-seeded EMA20/EMA50, strict two-candle-side fractals, pivot and confirmation
trend ordering, strict pivot placement, confirmation after t+2 and next-open
entry are fixed. EMA50 at confirmation fixes the stop; the target is 1.5R.
There is one position, no pyramiding, and invalid risk is rejected. The model
uses 5 bp taker fees and 2 bp adverse slippage per side, adverse gap stops,
no target price improvement, and stop first on ambiguous bars. Actual recorded
funding timestamps are retained, with conservative exposure when candle timing
cannot order the exit and funding. The exact simultaneous-open convention and
acceptable 0–1000 ms settlement offset are specified in the candidate. These
are simulation conventions, not verified exchange fills.

The estimator is total completed-trade net R divided by completed trade count.
Form daily pairs (S, N) and five consecutive 73-day batches. With theta as
that ratio and Z_j the batch sum of S - theta*N, the standard error is
sqrt((5/4)*sum(Z_j^2))/sum(N); the two-sided interval is theta +/-
2.7764451051977987 times that standard error. No bootstrap or alternative
method is selected on holdout outcomes. At least 200 completed trades, 30 exit
days, one trade in each batch, complete data and no open position are required.
For an adequate valid run a strictly positive lower bound means MET; a strictly
negative upper bound supports negative modeled expectancy; all other intervals
are inconclusive. Inadequacy produces INCOMPLETE. Invalid sources, missing
funding, incompatible schema, checksum mismatch or invalid inference variance
produce TECHNICAL_FAILURE where a run can start; missing source bytes may
instead prevent run reservation. A positive total return alone is not success.

Synthetic coverage tests support this method only for the prespecified
processes in `inference-validation.md`, not arbitrary long memory, regime
breaks or nonstationarity. The development calculation used 360 days and five
72-day batches, not the 365-day holdout method. The earlier failed seven-day
bootstrap and its preliminary intervals are preserved in the research record.

## Lifecycle and verification

The study-owned CLI uses the existing platform manifest, lineage, BeginRun,
source binding, canonical result publication and terminal seal. Its verifier
separately checks platform integrity and recomputes the complete trade ledger
and inference with the independent study replay. The replay and primary path
share the descriptive-summary helper: summary values have focused synthetic
tests, but are not claimed independently implemented. The implementation pins
the candidate spec and fails closed on changed source bytes, malformed input,
noncanonical results and semantic mismatch. A missing source may prevent a
run from beginning rather than producing a sealed technical-failure record.

Isolated synthetic workspaces, not the real study, exercised the platform
freeze/run/seal path. The synthetic manifest and corresponding test-only spec
digest differed from the actual candidate. Synthetic run
`20261029T093820Z-4ef367269ee3` completed with 1,477 simulated trades and an
inconclusive outcome; platform integrity and independent ledger/inference
replay passed. Separate synthetic run `20261029T093508Z-e5aa52b1096d` was
INCOMPLETE for an open final position and verified. A changed result artifact
failed integrity verification, and a missing source prevented run reservation.
A read-only final reviewer independently checked the synthetic COMPLETE seal,
the actual candidate hash guard, and the draft platform check. The actual study
has no frozen lineage or scientific run. These isolated test runs contain no
actual 2025 observations; their temporary harness is not a checked-in fixture.

Correction to the historical development audit: its reported single gap stop
was a floating-point reversal artifact in descriptive counting. Direct
comparison against the modeled slipped stop shows zero gap stops. This does
not change the trade ledger, net returns, inference, or scientific strategy.

## Rights and operator decision

The operator authorized noncommercial research using Binance Vision under the
applicable documented CC BY-NC-SA terms; this is not individual Binance
approval or endorsement. Do not redistribute raw provider data. Before a public
article, separately review attribution, share-alike treatment of derived
artifacts, noncommercial use and outstanding REST funding-data reuse rights;
a noncommercial site alone does not resolve every reuse question. Source terms
and evidence are recorded in `data-decision.md` and `development-evidence.md`.

Review and explicitly approve or reject the complete candidate before invoking
the platform freeze. Approval is a separate decision from authorization to
acquire and evaluate reserved holdout data. Freeze and holdout acquisition are
not performed in this milestone.

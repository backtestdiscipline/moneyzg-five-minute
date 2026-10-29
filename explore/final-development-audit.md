# Development semantic and statistical audit — exploratory

Classification: STATISTICAL_VALIDATION_FAILED. No freeze, holdout acquisition,
or 2025 result occurred. A failed synthetic calibration is a methodological
finding, not a holdout strategy outcome. Study and publisher evidence remain
under noncommercial restrictions; publication requires artifact-level review.

## Independent semantic reconciliation

A separate Python reconstruction, `explore/independent-replay/replay.py`, reads
the hash-verified development inputs independently, builds EMA and fractal
signals from their completed-bar windows, and then independently reconstructs
entry, position, costs, intrabar exits, gaps, funding and terminal state. Run
`go run ./cmd/development -ledger /tmp/moneyzg-dev-ledger.json` and
`python3 explore/independent-replay/replay.py /tmp/moneyzg-dev-ledger.json`
from the study root with Go 1.27.2 and offline module resolution. This is a
second implementation, not a platform integrity check or proof that prices were
executable. It directly compared all 1,959 closed trades, every recorded field,
and aggregate counts: 1,825 overlap signals ignored, 0 invalid entries, 4
stop/target ambiguities and no terminal position. There were zero mismatches and
zero floating-point differences for the recorded float fields. This verifies the
specified two-bar confirmation, t+3 entry, EMA band and trend at pivot and
confirmation, stop, 1.5R target, one-position logic, costs, gap rules and modeled
funding over the development paths. Intrabar exit timestamps identify the bar,
not an actual fill time; real funding settlement remains uncertain.

Independent review found a genuine defect: when more than one funding event
occurred during a trade, the funding-ambiguity flag could be overwritten by a
later unambiguous event. It did not affect P&L and there were zero ambiguous
funding trades in development. The flag now accumulates, with a synthetic
multi-settlement regression test; the full development ledger still matches.

## Development decomposition

1,959 closed trades: 1,057 long, 902 short; 719 net winners (36.7024%). Gross
mean +0.01922609 R; net mean -0.57728261 R. Per trade, modeled fees average
0.42538595 R, slippage 0.17015438 R, and signed funding cost 0.00096837 R.
The difference is dominated by the fixed assumed transaction costs, relative to
a small gross edge; this does not establish a live gross edge or profitable
execution at lower costs. There were 1,227 stop exits, 732 target exits, one
opening-gap stop (included in stops), zero terminal exits, four intrabar
stop/target ambiguous trades (0.2042%), and zero funding-ambiguous trades.
Initial-risk distance as basis points of slipped entry: nearest-rank p10
12.00311, median 31.50293, p90 76.55230. Modeled fees+slippage+signed funding
exceeded 0.5 R for 859 trades and exceeded 1 R for 271. These are fixed
assumptions, not observed fills; no tuning followed the diagnostic.

## Dependence and statistical calibration

Two separate checked-in synthetic simulations use 400 independent 365-day
series per case, stationary standard normal independent days or AR(1) with
phi 0.5 and 0.8 (innovation SD sqrt(1-phi²)), true daily mean zero. A fourth
case uses phi 0.8 with independent uniform daily trade counts in {0,...,4};
daily net R sums equal count times the daily shock, so the true trade-weighted
mean remains zero and zero-trade days are retained. Randomness is fixed
SplitMix64 with Box-Muller normals. See `explore/inference-primary/main.go`
for 10,000 bootstrap resamples with the proposed fixed 7-day circular blocks
and percentile ranks 250/9750. Coverage of zero: independent 378/400 (94.5%),
AR(0.5) 369/400 (92.25%), AR(0.8) 324/400 (81.0%), and AR(0.8) with varying
frequency 330/400 (82.5%). Monte Carlo standard error near 95% is ~1.1
percentage points; simulation is illustrative, not a model of the holdout.

`explore/inference-simulation/main.go` uses the same scenario design and 2,000
resamples per interval to compare blocks 3, 7, 14 and 28 days without choosing
from development performance. For AR(0.8) constant-count coverage was
68.25%, 81.25%, 89.25%, 90.25% respectively. This does not establish any of
those blocks as generally valid. The seven-day choice was a prior convention;
these checks do not justify it for potentially persistent dependence. Synthetic
coverage does not establish the actual holdout process and does not imply the
strategy is false. The preliminary development interval is not independently
validated as nominal 95% coverage.

Fixed descriptive 2024 block sensitivities (not candidate selection): L3
[-0.64026127,-0.51318094], L7 [-0.63676846,-0.51860417], L14
[-0.63438344,-0.51939875], L28 [-0.63827517,-0.51985828]. They do not
establish any interval's nominal coverage. Daily joint resampling preserves
within-block chronology and variable daily counts but necessarily breaks
long-range dependence across block boundaries. A pre-freeze statistical
redesign and independent calibration are required, not an opportunistic choice
of the most favorable interval or a lowered adequacy criterion.

## Platform and freeze gate

`records/moneyzg-five-minute/` is an unfrozen platform-generated draft. Its
manifest binds the development files and sidecars and declares post-freeze
2025 monthly candles, the two pre-holdout daily warm-up archives, their
publisher checksums and two bounded funding API responses. These paths remain
absent. The draft spec names both the unvalidated inference and unfinished
platform-bound holdout evaluation/semantic verification as blocking
uncertainties. `research-platform check moneyzg-five-minute` correctly reports
both blockers. No `freeze.json`, run or seal exists. Integration is partial:
the development evaluator is not yet wired to `BeginRun`, `OpenSource`,
immutable artifact publication, sealing or post-seal semantic replay. Those
study-owned integrations must be implemented and tested before freeze review;
platform integrity and scientific semantics remain separate claims.

The existing proposed 365 fully covered UTC holdout days, >=200 completed
trades, >=30 exit days, no open position and verified funding remain unchanged.
Before freeze, resolve a defensible dependence-aware acceptance procedure, test
its coverage and precision on independently justified scenarios, implement the
platform run/verification path, and review source licensing, contract and
funding event handling. Obtain explicit operator approval before freeze and
holdout acquisition. A source archive can be present on the server without
its observations being inspected; the holdout is historical, not a
calendar-forward prospective period.

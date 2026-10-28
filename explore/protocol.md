# Proposed confirmation-first protocol — implementation review only

This is ONE fixed interpretation for independent testing, not the creator's exact
method, a frozen record, or a performance result. The user's conventions resolve
source discretion: strict pivot price between EMAs, trend at pivot and confirmation,
fixed 1.5R target, fixed stop, one position and no pyramiding. No parameter search.

## Market and sources

Candidate: Binance USD-M BTCUSDT linear perpetual, last-trade five-minute candles,
quote currency USDT. The source video does not establish this exact contract. A
shortable linear contract supports both directions. Use publisher archives with
provenance, SHA-256 identities and publisher checksums where supplied; bind files
in the platform manifest. Require time-ordered timestamped trades for candle
construction, best bid/ask quotes with event timestamps and sequence identifiers,
and funding settlement rates; document exact source endpoints, coverage, license,
completeness and synchronization before freeze. The proposed method cannot be
executed credibly with OHLCV alone. Quote data feasibility and historical fee
schedule are unverified: if unavailable, stop, not substitute invented prices.

Development `[2025-01-01, 2026-01-01)` UTC is for implementation and debugging,
never parameter or threshold selection. Holdout `[2027-01-01, 2028-01-01)` UTC
is prospective and untouched, selected before observations; holdout bytes must
remain absent until after freeze and be declared post_freeze. Freeze must precede
both first holdout acquisition and holdout start; otherwise select a new future
interval by calendar rule, without examining outcomes, in a new draft. Acquire no
holdout now. Warm-up uses the 500 preceding complete bars separately for each
partition, without trades before the boundary. Use final confirmed market events
only; fail on missing required source, incomplete coverage, timestamp disorder,
invalid prices, impossible quote or publisher checksum mismatch. Events sharing a
timestamp need a publisher sequence; absent deterministic order means failure.
No event after the partition may enter an indicator, signal or fill.

## Signal and order semantics

Bars are UTC intervals `[open, open+5m)` aggregated from trades in event order;
open is first trade, close last, high/low extremes. Empty bars are data failures,
not filled forward. EMA at a bar close uses that close, alpha `2/(N+1)`, seeded
with first warm-up close, recursively in IEEE 754 binary64; periods 20 and 50,
no offset. Do not emit signals until 500 warm-up bars are complete. Preserve
bar closes and EMA values; no look-ahead. This seed and strict comparisons are
study conventions, not proven TradingView implementation. Never round indicator
values for comparison; validate finite values. Exact tie to any of four
neighbors disqualifies the strict high or low pivot; simultaneous high and low
pivot is ignored (ambiguous). The pivot's high/low is strictly between EMA20 and
EMA50 at the pivot close. For long, EMA20 > EMA50 at pivot and confirmation;
for short EMA20 < EMA50 at both. A pivot at t is determined only at close of
bar t+2, with the two neighbors on each side finalized. A signal with either
pivot or confirmation outside the partition is ignored. No future information
is used. Missing or revised data require a new byte identity and lineage.

An eligible signal submits a one-unit market order only after t+2 is finalized;
earliest eligibility is the following candle open boundary. To avoid claiming
that an adjacent candle's first trade was executable at the close, model one
second processing latency from that boundary and fill at the first recorded valid
best ask (buy) or bid (sell) at or after that time, within that candle. Quotes
must have a recorded event at or after eligibility; never reuse a stale prior
quote or fabricate an execution price. No quote in the candle: cancel the entry,
record an unfilled signal, and no later fill. This is a top-of-book price proxy,
not evidence of executable size, queue priority, or a real fill. Actual depth,
latency and impact must be validated; one-unit liquidity is a material blocker.

At fill, stop is EMA50 at confirmation, held fixed. For long require fill > stop;
for short fill < stop. Nonpositive or nonfinite initial price risk cancels the
trade before recording a fill, and is counted separately; no divide-by-zero.
The target is entry + 1.5*(entry-stop) for long, entry - 1.5*(stop-entry) for
short. This calculated level is an order trigger, never a synthetic fill. A
nonpositive target cancels the candidate. Exactly one open position; signals
while an order is pending or a position is open are ignored, with no pyramiding,
re-entry or reversal. A signal completing on a bar with an open position is
ignored even if an exit occurred earlier in that bar. A new signal may be used
only after the next complete bar. Stop and target are active immediately after
entry; use sequenced quotes strictly after the entry event. For a long, bid at or
below stop triggers a stop market sell, bid at or above target a target exit;
reverse using ask for a short. Fill each exit at the first quote on the required
side at or after trigger plus one second, within the partition; record the actual
observed quote, so gaps and adverse movement are possible. The model uses market
orders for both exits; target can miss or overshoot. A quote cannot trigger both
for valid stop/target levels. With an exit pending, ignore further triggers; if
no qualifying exit quote before partition end, mark open position INCOMPLETE.
Quote ordering resolves intrabar conflicts; if only OHLC bars exist, do not claim
an executable result or silently assume a stop/target order. Evaluate exits on
all quotes including the entry candle after entry. A last-bar unfilled entry is
cancelled at its boundary; a remaining position or pending exit is INCOMPLETE,
not silently dropped or artificially liquidated. Signals without room to confirm
and submit within the partition are ignored, with counts reported.

## Net accounting and inference

Proposed baseline fee: 0.05% of each entry and exit notional, charged as a cash
cost, not as a fabricated execution price. This is an explicit conservative
scenario, not a verified historical venue rate; verify it against venue fee
records and depth before freeze. Bid/ask incorporates observed spread; one-second
latency models delay but not unobserved market impact. Funding is charged at each
actual UTC settlement during the position, signed per official rate and mark
price (long pays positive rate, short receives); missing funding or mark is a
technical failure. One BTC notional is assumed as a normalized analytical unit;
no leverage, margin, liquidation or compounding claim. Before freeze demonstrate
available depth and capital/margin adequacy for that size or this execution model
is unsupported. P&L per completed position = signed exit minus entry, less entry
and exit fees and signed funding; divide by absolute entry-to-stop price distance
for net R. Risk distance excludes fees and funding; adverse gaps may exceed 1R.

Fixed sensitivity scenarios on the same signals and observed quotes: 0.04%,
0.05%, 0.10% fee per side, with the baseline 0.05%; also 5-second latency in
place of one second, evaluated separately and descriptively. These are scenario
assumptions, not independently verified fee or impact estimates. They are not
used to select the primary outcome. Gaps, cancelled entries, exposure, and all
cost components must be reported. No optimization and no use of a noncausal
marker strategy as the primary or a fair executable benchmark.

Primary metric: mean net R per completed trade. Use a deterministic circular
moving-block bootstrap of **UTC daily** net R sums and daily trade counts, jointly,
including zero-trade days, fixed 7-day blocks, 10,000 resamples, Go-specified
SplitMix64 seeded with 0x6d6f6e65797a6701, ratio of sums to counts. Generate
unbiased block starts by rejecting unsigned 64-bit draws below `(-365 mod
2^64) mod 365`, then taking draw modulo 365. Draw 53 independent block starts
per resample, wrap around the 365 days, and take the first 365 daily records.
Sort the 10,000 ratios ascending; use one-based nearest ranks 250 and 9750 for
the 2.5th and 97.5th percentiles. This estimates
uncertainty conditional on the selected market and period; seven days is a fixed
assumption, not proof of adequate dependence coverage. Use IEEE 754 binary64 for daily sums and ratios, summed chronologically in
each resample; verify reproducibility against checked-in synthetic fixtures. Adequacy requires 365 fully covered holdout UTC days, at least 200
completed trades, at least 30 distinct days with trades, nonzero bootstrap
resample counts and no unresolved position at the boundary. These pragmatic floors
are not a power guarantee; assess precision/power on synthetic data before freeze.
If coverage or count floors fail, scientific result INCOMPLETE; invalid or
corrupt inputs, failed invariants, unsupported source fidelity, or replay mismatch
are TECHNICAL_FAILURE. A valid but wide interval is COMPLETE and inconclusive.
A positive expectancy claim is accepted only if the 95% interval's lower bound
is strictly above zero and data adequacy holds. A 95% upper bound strictly below
zero supports negative expectancy; all other valid intervals are inconclusive.
A positive aggregate return, alone, cannot pass. Predefine sensitivity as
descriptive, not multiple confirmatory tests. Study-owned semantic verification
must recompute decisions, cost accounting, outcome and inference; replay with the
recorded committed implementation and independently review the algorithm.
Platform integrity PASS is distinct from semantic or economic validity.

## Gate

Ready for implementation review of this draft only. Freeze review requires:
identified licensed quote/trade/funding sources and checksum coverage; temporal
alignment, liquidity/depth and fee evidence; a deterministic implementation,
semantic tests and independent review; inference precision justification; manifest and prospective holdout definitions; and a
fully specified frozen acceptance rule. Holdout stays absent. If a data source
cannot support observed quote fills, do not relabel an OHLC approximation as real
execution. No final evaluation or result is authorized by this protocol.

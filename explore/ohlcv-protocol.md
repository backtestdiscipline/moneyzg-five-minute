# OHLCV evaluation and freeze proposal — unfrozen

## Scope and source

On operator instruction, the previously unfrozen quote-level proposal in
`protocol.md` is superseded by this lower-fidelity five-minute-candle model.
The archive-size, liquidity, data quality and rights findings remain in
`data-feasibility.md`. This avoids the ~63 GB quote/trade requirement and does
not represent actual execution. This is a causal independent interpretation,
not an exact reproduction or a production trading system. No strategy parameter
was selected using returns. There is no 2R sensitivity in this proposal.

Source evidence and timestamps remain in `reconstruction.md`: the video shows
20/50 close EMAs, a two-period fractal, the EMA band and discretionary exit
choices. The operator's clarified interpretation and the choices below specify
a single mechanical version; neither 1.5R nor a fixed stop is a mandatory source
rule. The retrospective examples do not establish look-ahead bias. TradingView's
own description identifies a five-candle pivot and a two-or-more-candle lag:
https://www.tradingview.com/support/solutions/43000591663-williams-fractal-indicator/.
The specific implementation's equality and historical data revision semantics
are not established; strict ties below are our convention, not a claim of exact
TradingView equivalence.

## Data and split — conditional, not approved

Candidate instrument is Binance USD-M BTCUSDT linear perpetual, last-trade 5m
candles, UTC. A perpetual supports both directions; fixed one BTC units are
used only to normalize per-trade P&L, with no leverage, liquidation, collateral
or account-level feasibility claim. Historical contract rules and fees must be
verified for the period. Funding is necessary: require each 00:00, 08:00, 16:00
UTC settlement rate and contemporaneous settlement mark, charging positive
funding to longs and crediting shorts. If settlement mark or funding interval
cannot be verified, fail technically rather than use zero or a later mark.
Use official publisher archives and/or documented funding-history API, with
bytes acquired and hashed before parsing, full source locators, dates, checksums,
permissions and platform manifest binding. Binance monthly 5m archives exist
for May 2023–March 2024 per metadata in the feasibility report; a May 2023 file
was independently checksum verified. Later candidate-period availability,
quality, and exact settlement mark access remain unverified. Prior BTC study's
spot 1h sources are not the same instrument or timeframe. The available May
sample lacks a verified 500-bar warm-up and settlement marks; no development
strategy diagnostics were run.

Candidate development `[2024-01-01T00:00Z, 2024-12-30T06:20Z)` and
holdout `[2025-01-01T00:00Z, 2026-01-01T00:00Z)` are selected without
reviewing strategy returns. The 500-bar holdout warm-up begins 2024-12-30
06:20 UTC and is reserved, unacquired and excluded from development; the
development warm-up begins 2023-12-30 06:20 UTC. Reserve holdout and its
warm-up as post-freeze sources; acquire and bind them only after freeze.
This is a completed historical interval reserved prospectively to the freeze,
not a calendar-prospective test. Public archive HEADs returned 200 for all 2024
and 2025 monthly candle files but do not prove quality or complete coverage.
The change from the previous 2027 proposal is documented in `data-decision.md`.
No missing, duplicate or invalid bars are filled. Publisher sidecars and local
hashes prove byte identity, not correctness of market data.

The Binance Vision dataset terms describe generally CC BY-NC-SA 4.0, noncommercial
research permissions and a separate license for commercial use:
https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md.
A potentially commercial Backtest Discipline article and derived results need a
rights determination before using this source for a publishable study. No claim
is made that publication is permitted or prohibited in this particular case.
An alternative source is not silently selected.

## Rules implemented in `ohlcv/evaluator.go`

Inputs are exactly 500 contiguous 5m warm-up bars followed by complete contiguous
bars in `[start,end)`, plus all 8-hour settlement rate/mark records inside that
interval. Bar timestamps are UTC milliseconds at open; prices must be finite and
positive, high/low contain open and close, volume finite and nonnegative. No
interpolation. Funding records are unique, aligned, finite and positive in mark;
a missing or invalid input returns an error. Upstream ingestion must verify
source bytes, close-time convention, warm-up, final bar coverage, data rights,
settlement intervals and the contract. The evaluator accepts typed records; it is
not yet a decoder, platform-integrated run, or semantic replay implementation.

Seed both EMAs with the first warm-up close and update recursively using
`alpha=2/(N+1)` in binary64, N=20 or 50. At each completed candle i, consider
pivot p=i-2 only when p is in the evaluation partition. A strict low must be
below both preceding and both subsequent lows; strict high is the reverse.
Ties exclude a pivot, as does simultaneous high and low. The low must be strictly
between EMA50 and EMA20 at p and EMA20>EMA50 at p and i for a long; the high
must be strictly between EMA20 and EMA50 with EMA20<EMA50 at p and i for a
short. Values at p are fixed at its close; no future data are used. The signal
exists at close(i) and, if flat, enters only at open(i+1), the pivot's t+3.
Signals overlapping an existing position are ignored. An exit within a bar can
be followed by a newly confirmed signal at that bar's close, entering at the
next open; no reversal at the exit event. No pyramiding.

Entry is a simulated open-price proxy plus adverse 2 basis points for a buy or
minus 2 bp for a sell. This slippage is an explicit stress scenario including
unobserved spread/impact, not a verified bound or actual fill; open may have the
same nominal timestamp as confirmation close, and delay is not modeled. Initial
stop is fixed at confirmation EMA50. Require strictly positive entry-to-stop
risk; otherwise skip and count. Target is 1.5 times this risk beyond the
slipped entry; a nonpositive target is skipped. One BTC is modeled per position.

Stops and targets are active on entry candle. On each candle, evaluate opening
price first: a stop gap fills at the worse observed open; a target gap fills at
the target with no favorable gap improvement. If the open is inside the levels
and both high/low cross stop and target, count one ambiguous trade and choose
stop first. Otherwise fill at the touched stop or target level. Every exit is
then shifted adversely by 2 bp. These are price proxies, not observed actual
fills; gaps within a candle may be worse than the model. Charge taker fee 5 bp
of slipped notional on each entry and exit. Fees and slippage are predeclared
scenarios, not proven historical rates or conservative bounds. Funding at an
8-hour candle open is charged to positions already open before exit or entry at
that boundary; no funding on a new position at that boundary. Signed funding is
rate × settlement mark × direction × one BTC. Binance documents deviations in actual settlement processing, and a small
2024 API sample contained timestamps 1 ms after the nominal boundary. The
current evaluator rejects those records: a reviewed mapping and boundary policy
is required before integration; see `data-decision.md`.

Gross R is direction × (raw exit proxy minus raw open proxy) divided by initial
risk; net R uses slipped entry/exit, minus both fees and signed funding, divided
by the same risk. Return per-trade components and count invalid entries, ignored
signals and ambiguous trades. An open position at end produces INCOMPLETE: it is
not liquidated, silently omitted or carried into a holdout; no pending entry is
created at the final bar. Avoid interpreting partial completed trades as the
primary result. The number and percentage of ambiguous trades are reported over
completed trades, and zero trades has an undefined percentage. Open/close
availability and all intra-bar execution are idealized; OHLC cannot resolve
spread, queue, impact, intrabar gaps, nor actual stop-versus-target chronology.

## Hypothesis, adequacy and deterministic inference

Primary question: whether this fixed causal interpretation has positive mean net
R after the stated modeled transaction costs and funding. Report trade count,
win rate, gross and net expectancy, profit factor (positive net R sum divided by
absolute negative net R sum; undefined if no losses), aggregate fees, slippage
and funding, ambiguous count/percentage, and daily series. Report maximum
peak-to-trough drawdown for a hypothetical **fixed initial-risk unit** convention:
initial wealth 100 R, add each closed trade's net R, no compounding, ordered by
exit bar. Divide peak-to-trough loss by the preceding peak wealth; if wealth is
nonpositive report the event and stop interpreting percent drawdown. This is not
an exchange account, leverage or liquidation simulation.

Keep the former threshold: 365 fully covered holdout UTC days, >=200 completed
trades, >=30 distinct UTC exit days, no open position, and nonzero counts in every
bootstrap resample. Inputs with missing settlement records or invalid bars are
TECHNICAL_FAILURE; inadequate count/coverage is INCOMPLETE, not negative evidence.
For valid complete data, calculate mean net R and a two-sided 95% circular
moving-block bootstrap confidence interval for the ratio of sums of daily net R
to sums of daily trade counts, including zero-trade days. Resample 365 days in
7-day circular blocks, 53 blocks per resample, 10,000 resamples, fixed SplitMix64
seed `0x6d6f6e65797a6701`, rejection sampling of unsigned 64-bit values below
`2^64 mod 365`, modulo 365 for block starts. Use first 365 records, binary64
sums in drawn order, sorted one-based ranks 250 and 9750. Assign trades to exit
UTC day. Lower interval bound >0 is SUCCESS for this modeled hypothesis; upper
bound <0 is FAILURE to find positive expectancy with negative estimate supported;
otherwise COMPLETE and INCONCLUSIVE. Positive total return is insufficient.
These pragmatic sample floors are not a power guarantee. Test bootstrap
implementation with synthetic data and assess coverage/precision under plausible
dependence before freeze; 7-day blocks are an assumption, not proven adequate.
Robustness findings should not be selected or tuned from holdout results.

## Freeze review blockers

The signal/execution core and synthetic tests are exploratory, not a complete
study: source permissions and historical contract/fee/funding semantics must be
resolved; data acquisition/decoder and publisher checksum validation, platform
manifest and post-freeze bindings, reporting, inference, deterministic semantic
verification and independent review are still needed. Verify interval coverage
and full source rights, and independently justify cost assumptions and inference
precision. Only then prepare and review platform draft records and final
acceptance rule. No freeze or holdout evaluation is authorized now. A valid
negative result concerns this model, not a proof the creator's strategy is false.

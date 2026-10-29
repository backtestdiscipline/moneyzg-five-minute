# MoneyZG five-minute strategy: publication evidence handoff

**Editorial preparation only.** Empirical tables, figures and claims in this
package remain conditional on the rights review in `licensing.md`; the package
itself does not grant publication permission. No raw market data, original
video or transcript, complete ledger, or internal platform records are included.

## Source and reconstruction

MoneyZG, [“Easy 5 Minute Scalping Strategy (Simple to Follow Scalping Trading
Strategy)”](https://www.youtube.com/watch?v=nFYE0H-Lfzg), 14 May 2021. A
supplied local video and transcript were consistent with the indexed identity;
their original acquisition and unaltered publisher provenance were not
independently established. Timestamps are approximate video offsets:

- 00:05–01:09, 01:48–02:11 and 03:25: five-minute chart, multiple markets;
  the displayed chart says BTCUSDT. The source does not establish the exact
  contract or price feed used in our study.
- 02:31–03:40: EMA20 and EMA50, close source, offset zero.
- 03:51–04:13 and 06:27–07:59, settings at 06:57: Williams Fractals, Periods
  2; the creator says to wait for two further candles to complete.
- 08:08–09:24: trend ordering and low/high fractal between the EMAs for
  long/short setups; stop at EMA50. At 10:19–10:54 and 11:54–12:10 targets
  are discussed as risk multiples with discretion, not a mandatory fixed 1.5R.
- 14:32–15:33: adding positions is optional. At 18:37–18:48 the creator warns
  about fees; no numerical fee, historical win rate or measured return was
  established in the inspected material.

The illustrated trades are retrospective. They do not establish a recommendation
to trade unconfirmed fractals. The original leaves tie handling, exact EMA
comparison time, order type, fill, stop adjustment, target and discretionary
filters insufficiently defined for a unique executable reproduction.

## Fixed study interpretation

The study used Binance BTCUSDT USD-M perpetual five-minute candles and funding
records. EMA20 and EMA50 use close, seeded from the first warm-up close. A
strict five-candle low/high pivot must lie strictly between its pivot-time EMAs;
trend ordering must hold at the pivot and confirmation. A pivot at t becomes
known only after t+2 closes. The earliest modeled entry is the t+3 candle open;
the retrospective marker belongs to t and is not a tradable t-time signal.
Strict ties are a study convention, not verified as identical to TradingView.

The stop is fixed at confirmation-time EMA50, the target at 1.5 initial R, and
one position is allowed with no pyramiding. The video permits discretionary
choices; these mechanical rules are the study's, not verbatim source rules.
Entry and exit use 5 bp taker fees and 2 bp adverse slippage per side. Candle
high/low determines contact, gaps use adverse stop prices, target fills receive
no price improvement, and ambiguous stop/target bars use stop-first. Funding
uses recorded timestamps and marks, with conservative exposure when candle
resolution cannot order an exit and settlement. Entry opens are price proxies,
not verified live fills. R is each trade's initial slipped-entry-to-stop risk.

## Results and uncertainty

The predeclared question was whether the fixed interpretation has positive
mean net R per completed trade after modeled costs. A positive lower bound of a
95% interval, plus all data and sample adequacy requirements, was required for
acceptance. The untouched historical holdout was 2025-01-01 through
2025-12-31 UTC; 500 prior bars warmed indicators without eligible warm-up
pivots. The holdout was reserved before the freeze but is not calendar-forward
validation. Development in 2024 was exploratory: 1,959 trades, mean net
−0.577283 R; do not pool it with holdout results.

The frozen holdout was **COMPLETE**, with **NEGATIVE_EXPECTANCY_SUPPORTED**:
1,833 completed trades, mean gross +0.008310 R and mean net −0.663233 R; 95%
interval [−0.783204, −0.543262] R. Net win rate was 35.19%. There were 938
longs, 895 shorts, 1,164 stops and 669 targets. Three trades had ambiguous
stop/target candles; ten had ambiguous funding exposure. Mean modeled fees,
slippage and funding were 0.479361, 0.191744 and 0.000438 R per trade,
respectively: costs account for the shift from slightly positive gross to
negative net expectancy. Slippage is included in net fills and must not be
deducted a second time. These are modeled costs, not observed account debits.

The interval uses five consecutive 73-day batches of paired daily net-R sums
and completed-trade counts, including zero-trade days, a ratio-influence
standard error and Student-t with four degrees of freedom. There were 355 exit
days and each batch had trades (383, 371, 359, 346, 374); the frozen minimums
of 200 trades, 30 exit days, complete data and no open final position were met.
Synthetic coverage validation is conditional on the processes tested and does
not guarantee coverage under long memory or regime changes.

Platform integrity and independent signal, ledger, funding and inference replay
passed. The descriptive summary helper was shared and is not independently
implemented. Five-minute candles do not reveal actual intrabar order, spread,
queue position, impact, fills or precise funding eligibility. The findings do
not disprove all interpretations of the original strategy. `summary.json` holds
precise values and evidence hashes; `chart-data.csv` provides aggregate units.
The sealed authoritative result and verification remain in the study repository.

# Frozen MoneyZG 2025 holdout result

Run `20261029T132412Z-59e969df48e1` is **COMPLETE** and its frozen scientific
classification is **NEGATIVE_EXPECTANCY_SUPPORTED**. Mean modeled net R per
completed trade was -0.663233, with a two-sided 95% five-batch ratio-influence
Student-t interval of [-0.783204, -0.543262]. Its upper bound is below zero;
the criterion for a positive expectancy claim, a lower bound above zero, was
not met. This is evidence about the specified candle-based model and costs, not
proof of actual executable returns or a universal falsification of the video.

The study used the reserved historical 2025 holdout, not a calendar-forward
experiment. Prior 2024 development diagnostics were disclosed before freeze;
no model or acceptance criteria were changed after seeing holdout results.
The full-development diagnostic had 1,959 trades and -0.5773 mean net R; its
related 360-day, five 72-day-batch interval was [-0.669960, -0.489376]. It
is exploratory and is not a 365-day holdout interval.

## Results and adequacy

- 1,833 completed trades: 938 long, 895 short; 645 net winners (35.19%).
- Gross expectancy +0.008310 R; net expectancy -0.663233 R per trade.
- Total modeled fees 185,701.28616451, adverse slippage 74,280.51524689 and
  net funding charge 38.35497200 USDT for the mechanical one-BTC-per-trade
  convention. These are simulated, not account debits or account returns.
  Dividing each trade's costs by its initial risk gives total 878.668058 R in
  fees, 351.467224 R in slippage and 0.803127 R in funding, respectively.
  Slippage is already reflected in filled prices and net R; do not subtract it
  again from net R. Gross expectancy is the unslipped, before-cost comparison.
- 1,164 stop exits, 669 target exits, zero identified gap stops, zero terminal
  liquidations. Three trades (0.164%) had both stop and target contacted in a
  candle and used stop-first; ten had ambiguous funding exposure, handled by
  the frozen conservative convention. The run ended with no open position.
- 355 distinct UTC exit days; five chronological batch counts are 383, 371,
  359, 346 and 374. All frozen minimums pass: >=200 completed trades, >=30
  exit days, >=1 trade in every batch, complete valid inputs and no open final
  position. The five batches contain 73 days each, including zero-trade days.

## Sources and verification

The 14 Binance Vision archives and their publisher checksum sidecars plus two
public funding API responses were acquired between 2026-10-29T13:23:14Z and
13:23:49Z after freeze. `provenance/holdout-acquisition.json` records all 30
exact public URLs, individual timestamps, sizes and SHA-256 values; raw bytes
stay in ignored `data/`. All 14 publisher checksums matched. Independent
prechecking found 500 warm-up and 105,120 contiguous evaluation candles, valid
OHLC and close timestamps, and 1,095 distinct expected eight-hour funding events
with finite rates, positive marks and actual offsets 0–16 ms from their nominal
boundaries. The study decoder and evaluator also validate source identity,
schema, continuity and funding under the frozen rules.

The platform execution plan binds the frozen lineage, acquired source
identities and implementation. Result SHA-256 is
`4ead222258da57befb7f662230d24706c17177940f5484f0f81ced2c918eca84`;
seal SHA-256 is
`bb2297ac355a1c2727077fd9232a70cbabcc408fba90e854a22d2c075f94c463`;
specification SHA-256 is
`14d326cf92bb92fefd5ef65fcaa6fe8e074aeea75aca3f9073df86ad1b44697c`.
`research-platform verify` passed. The study `verify` reported matching
implementation, integrity PASS and semantic PASS, independently reconstructing
signal timing, the complete ledger, funding and five-batch inference. The
primary and independent replay share the descriptive summary helper; its
numbers are tested, but that helper is not independently implemented.

Candles cannot recover actual intrabar ordering, spread, queue position, market
impact or fills. Conservative assumptions cannot establish realized execution.
The interval's synthetic validation covers specified processes, not arbitrary
long memory or regime change. The Binance Vision noncommercial basis and
attribution and share-alike obligations are documented in the existing study
record. A later public article requires a separate rights and attribution
review, including applicable funding REST-data terms; do not publish raw data.

# Study charter — MoneyZG five-minute scalping

State: exploratory source reconstruction; no lineage frozen, no run, no result.
Kind: proposed original independent empirical study, motivated by a video; exact
reproduction and a causal adaptation are distinct possible targets.

Questions: RQ1 can the original rules be reconstructed? RQ2 are illustrated
entries causally possible? RQ3 does a fixed, causal implementation have positive
net expectancy on untouched data? RQ4 how sensitive is it to execution costs and
fills? RQ3 is conditional on a sufficiently specified causal method; an
adaptation must be explicitly identified and independently justified.

Non-goals: trading, parameter optimization, performance claims from selected
examples, final evaluation, holdout outcome inspection, or publication now.
Source permissions: public web pages are used only for identifying and
interpreting the source; video, captions and market data have not been acquired.
Permissions to retain or share source bytes and excerpts must be established
before acquisition or external AI sharing. Do not commit source bytes. Capture
and hash permitted original bytes before using them as evidence; record locator,
acquisition time, access method and provenance. Search snippets below are
unhashed secondary representations, not a preserved primary source.

## Prospective design outline — NOT a frozen charter

Source-derived method, exact instrument and contract, dates, signal, sizing,
exit and confirmation are unresolved. BTCUSDT perpetuals are a candidate because
shorts need a shortable instrument; funding, liquidation, mark/index versus last
price and exchange rules must be modelled. Spot would need a separately justified
shorting mechanism. Do not select either based on observed performance.

Before development outcomes, specify fixed UTC development and later untouched
holdout intervals, venue, file locators, publisher checksums, retention and
post-freeze acquisition, plus warm-up history and partition-boundary treatment.
No holdout outcome inspection before the final acceptance rule is frozen. Use
verified contiguous OHLCV, validated prices, timestamps and source identity;
assess spread, slippage, fees and funding against contemporaneous independent
sources, not arbitrary optimistic defaults. OHLCV alone cannot establish
executable bid/ask, intrabar ordering or liquidity. Prespecify adverse fill
semantics and a cost sensitivity grid; same-bar stop/target collisions need a
conservative rule or finer verified data. Record decision, order and fill times;
a close-derived decision cannot assume an achievable fill at the same instant.
Fix EMA seed and warm-up, equality rules, overlapping trades, pyramiding,
funding and price gaps. No optimization of source parameters or exits.

Candidate primary metric: mean net return per unit of initial risk per eligible
trade, with explicit risk denominator, aggregation and handling of open trades.
Specify a dependence-aware confidence interval and fixed block or HAC parameters
before freeze; zero or negative lower confidence bound does not establish positive
expectancy. Fix minimum sample, time coverage, power/precision rationale and
criteria for inadequate data before seeing holdout outcomes. A positive total
return alone is insufficient. Missing/corrupt sources, invalid data, unavailable
execution assumptions or broken replay are technical failures, distinct from a
valid negative or inconclusive result. A complete deterministic implementation,
synthetic causal and execution tests, independent semantic review, platform
integrity verification and replay must precede closure. These are design
requirements, not a completed or accepted statistical plan.

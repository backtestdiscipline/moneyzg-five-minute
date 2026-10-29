# Data rights, funding and holdout decision — 2026-10-29

Decision: **rights and funding-timestamp semantics block implementation toward
freeze review**. The OHLCV model remains authorized and unfrozen. No change to
strategy, costs or acceptance floors; no holdout data or performance inspected.

## Rights and dates

Binance [Vision Dataset Terms v1.0](https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md)
are labeled last updated 2026-08-26, not a separately verified effective date.
They describe free CC BY-NC-SA 4.0 access, including personal non-production
backtesting and noncommercial research, with attribution/share-alike for
redistributed derivatives; they say commercial use requires a separate written
enterprise data license. They govern the dataset and its derived outputs. Terms
say amendments apply prospectively to datasets made available after an amendment
and do not alter rights to data already downloaded; this does not establish the
applicable license for an old archive first downloaded today. The prior study
sample acquisition on 2026-10-28/29 is recorded in `data-feasibility.md`; its
local temp-file timestamps agree. Original publication of an archive in 2023 or
2024 does not establish an earlier acquisition by this study. No holdout archive
has been obtained. No historical terms or alternative permission has been proven.

Outstanding question: may Backtest Discipline, given its actual monetization and
distribution model, publish the specific derived empirical tables, charts and
article from this research, under which attribution and share-alike conditions,
and does a written commercial license apply? Seek a written rights determination
from Binance or a qualified rights holder for the actual use; do not infer it
from silence or the accessibility of the files. This is an evidence assessment,
not legal advice. The existing study's small noncommercial development samples
do not establish publication permission.

For keyless `fapi.binance.com` REST data, documentation establishes an API
interface, not a license to publish. Vision terms mention `data.binance.vision`
and “associated endpoints” but do not conclusively classify this endpoint or the
relevant jurisdiction and use. The [API Keys Product Terms](https://bin.bnbstatic.com/static/cms/cg08ou2ak0tn7mcplvfg/file/826cc99cc3847cd017dc7b82525896f01ac211d9c177573d8ea958402766d965.pdf),
effective 2026-01-05, address API keys; the sampled endpoint used no key, so
those terms alone cannot settle its rights. Identify the applicable Binance
entity, REST terms and written permission for derived public/commercial results;
REST access is not a proven workaround for Vision restrictions.

A practical alternative to investigate, subject to separate operator approval
before purchasing, is CoinAPI paid data: its [commercial-use policy](https://www.coinapi.io/blog/coinapi-data-commercial-use-policy)
explicitly describes commercial backtesting and sharing non-reconstructable
summaries and charts on paid plans, while prohibiting raw redistribution. Its
[metrics catalog](https://www.coinapi.io/datasets/metrics-ohlcv) describes
funding and mark metrics, but exact BTCUSDT perpetual 2024–25 coverage, precise
settlement marks, cost, plan terms and right to publish this particular report
must be confirmed in the applicable contract. It is not a verified free option,
no subscription was purchased, and provider policy is not a substitute for the
contract governing the intended use.

## Funding feasibility and observed timestamp mismatch

On 2026-10-29 approximately 07:26 UTC, one unauthenticated request to
`GET https://fapi.binance.com/fapi/v1/fundingRate?symbol=BTCUSDT&startTime=1704067200000&endTime=1704671999999&limit=100`
returned HTTP 200, 2,647 bytes, SHA-256
`6937e2c3980e2495e34d9dfc24f0b46385496376c77ec349ae0d4fa13854481e`.
The sample is `/tmp/moneyzg-funding-2024.json`, temporary and outside Git;
no publisher checksum is supplied by this API. Hash was calculated before
parsing. It contains 21 unique BTCUSDT regular records, covering Jan 1 00:00
through Jan 7 16:00:00.001 UTC; rates and positive marks are finite. Eighteen
funding times lie exactly on the 8h grid, three are 1 ms after it; adjacent
intervals were 28,800,000 ms (17), 28,800,001 ms (2), and 28,799,999 ms (1).
There is one event per nominal settlement in this week, not proof of full
2024 or 2025 coverage. No holdout values were requested.

[Official API documentation](https://developers.binance.com/docs/derivatives/usds-margined-futures/market-data/rest-api/Get-Funding-Rate-History)
describes funding time, rate and the mark associated with the funding fee, so
these records can provide rate and mark together and separate mark-price candles
are not inherently necessary, subject to complete validated records and rights.
Binance's [funding FAQ](https://www.binance.com/en/support/faq/detail/360033525031)
says funding is position notional at mark times rate, with a default 8h interval
that can change and a possible 15-second deviation in actual fee transaction
time. It does not prove how to allocate a hypothetical open or exit at that
boundary. This is a documented modeling limitation; it must not be treated as
exchange-exact execution.

The current evaluator requires exact 8h timestamps and rejects off-grid values;
a synthetic test now confirms rejection. Do not silently round these records or
change the funding charge order. Before freeze, specify a deterministic validated
mapping from API times to nominal settlements, collision/missing-event rules,
interval changes, and the allocation of exposure at the boundary, while retaining
raw timestamps for audit and explicitly qualifying the 15-second uncertainty.
Until then the actual sample is incompatible with the evaluator. Synthetic
long/short sign and fee-accounting tests pass; they are not an API or exchange
settlement verification.

## Proposed historical split and outstanding integration

Recommend 2024 development through 2024-12-30 06:20 UTC, then reserve 500 bars
for holdout warm-up and 2025 as the 365-day holdout. Precise half-open intervals
are in `ohlcv-protocol.md`. The split was chosen by calendar and clean warm-up
separation, not outcomes. It is prospective only relative to the future freeze,
not to the market dates. On 2026-10-29 public HEAD requests returned 200 with
nonzero content lengths for all 24 monthly 2024–2025 5m candle archives; total
bytes and mere availability do not establish continuity, checksums or rights.
No holdout archive or strategy output was retrieved. Retain 365 days, 200 trades,
30 trading days and the earlier confidence-interval acceptance criterion.

Once rights are established, minimal remaining work is verified source ingestion
and platform source binding, a reviewed funding mapping, quality checks, result
summaries and fixed bootstrap, independent semantic verification and freeze review.
Do not implement bootstrap or acquire full data while the rights question blocks
the intended study; no generalized provider or funding service is needed.

# Data and execution feasibility — 2026-10-29

**Classification: BLOCKED_BY_DATA.** Historical quote archives exist, but the
current proposed 2027 holdout cannot be assumed available, the 2023–24 candidate
is not validated for continuity or 1 BTC execution, and the short holdout
requires an explicit scientific protocol amendment. Do not implement the complete
backtester or freeze. This is a study review label, not a platform lifecycle state.

## Required inputs and optional improvements

The existing `protocol.md` requires individually timestamped trades to construct
five-minute last-trade bars, sequenced best bid/ask updates and quantities for
entry and exit proxies, funding rates and corresponding settlement mark price,
and a justified fee/latency model. Exact quote ordering and synchronization with
trade bars and funding events must be specified. A 500-bar warm-up is required.
Publisher 5m candles are a useful independent cross-check; they do not replace
trade input without amending the method. One-minute mark candles are a check, not
necessarily the mark at a funding instant. Full order-book replay, multi-venue
quotes and greater depth are optional improvements; they are not promised by
this gate. Best bid/ask quantities are essential to assess the proposed 1 BTC
size, but even adequate displayed size cannot guarantee a fill. Historical
contract filters, fees, order and settlement rules are necessary for an economic
claim and remain unverified for the proposed dates.

## Archive metadata and coverage

On 2026-10-29 (UTC session spanning 2026-10-28/29), public HTTPS HEAD requests
to `https://data.binance.vision/data/futures/um/monthly/` for BTCUSDT returned
200 with nonzero content lengths for **every month May 2023–March 2024** for
`bookTicker`, `trades`, `aggTrades`, `klines/5m`, `markPriceKlines/1m` and
`fundingRate`. This verifies object presence and byte lengths, not content
completeness, tick integrity or suitability. All these are publicly reachable
without a private credential; no holdout files or outcome values were downloaded.
BookTicker monthly size ranged from 1,981,788,217 bytes in partial May 2023 to
8,270,743,447 bytes in February 2024. April 2024 also returned 200 and
37,761,761 bytes: its existence must not be mistaken for a full month. A Binance
contributor confirmed discontinuation in May 2024 and the issue reports April's
file contained only early hours: https://github.com/binance/binance-public-data/issues/334.

The independently retrieved May monthly bookTicker and trades `.CHECKSUM`
sidecars returned 200. Binance says each zip has a SHA-256 sidecar:
https://github.com/binance/binance-public-data#checksum. These sidecars have
not been checked for every object. HTTP HEAD, ETags, and claimed dates do not
verify bytes. The separate development samples below passed published checksums;
no complete-period gap audit was performed. A user analysis reports gaps in
BTCUSDT trade archives, including September and November 2023; treat this as a
lead, not an independently verified result:
https://github.com/binance/binance-public-data/issues/424. A third-party derived
dataset reports bookTicker beginning 2023-05-16 11:49 UTC and coverage through
March 2024; its completeness claim is not independent verification of raw event
continuity: https://huggingface.co/datasets/Mindbyte-89/btcusdt_perp_bookticker_features_1m_05_2023_to_03_2024.

## Bounded development samples, schemas and verification

Four files from May 2023 only were downloaded into `/tmp/moneyzg-feasibility/`,
outside Git. The data have not been used to run signals or assess performance.
The listed hashes were calculated before reading their contents and matched the
publisher sidecars. These temporary copies are not durable archival retention.

| Archive | Compressed bytes | SHA-256 |
| --- | ---: | --- |
| `daily/bookTicker/BTCUSDT/BTCUSDT-bookTicker-2023-05-16.zip` | 53,816,320 | `1e51fcd332f9c901490a6e384dcd9b21943c04240180bda4903593a37bfc6262` |
| `monthly/klines/BTCUSDT/5m/BTCUSDT-5m-2023-05.zip` | 401,781 | `def79dfc73dc7031f961a7cdd83eb9f815a0775502829201e793876b61065eeb` |
| `monthly/markPriceKlines/BTCUSDT/1m/BTCUSDT-1m-2023-05.zip` | 1,007,316 | `fa5ad18ee67b48ffde94a57c674900a7129c1f5078a3d82c02f4cd8ec3411e42` |
| `monthly/fundingRate/BTCUSDT/BTCUSDT-fundingRate-2023-05.zip` | 983 | `5d026fece46df86293054c8c56c2c8de575d30161b82979b11831831217a246e` |

The daily quote CSV header is `update_id,best_bid_price,best_bid_qty,
best_ask_price,best_ask_qty,transaction_time,event_time`; timestamps in this
sample are Unix milliseconds. Its 4,494,018 rows span **11:49:47.209 to
21:44:00.522 UTC on May 16**, not a complete UTC day. Within these observed rows,
transaction times never decrease, update IDs always increase, 1,604,838 adjacent
pairs share a transaction timestamp, maximum consecutive gap is 1,372 ms, and
no nonpositive or crossed price was observed. This cannot prove no dropped
updates: IDs are not necessarily consecutive and no separate sequence guarantee
has been established. 1,917,914 rows have less than 1 BTC displayed on at least
one side; that is a row-weighted observation in a partial day, not a proportion
of time, future trade signals, or an estimate of impact. It makes assuming
unconditional 1 BTC at the best quote indefensible. The full-day absence after
21:44 may reflect an archive gap; no inference about subsequent days is made.

The sampled 5m and mark 1m CSVs have 12 fields: `open_time,open,high,low,close,
volume,close_time,quote_volume,count,taker_buy_volume,taker_buy_quote_volume,
ignore`. The funding CSV has `calc_time,funding_interval_hours,last_funding_rate`
and no settlement mark. Official USD-M funding-history API documentation describes
`fundingTime`, `fundingRate` and the mark associated with the charge:
https://developers.binance.com/docs/derivatives/usds-margined-futures/market-data/rest-api/Get-Funding-Rate-History.
Historical API coverage and a checksum-backed settlement-mark capture are not
verified. A one-minute mark candle alone does not prove the exact mark at
settlement. Individual and aggregate trade ZIPs were confirmed only by HEAD, not
sampled or schema-validated; public documentation describes trades and aggregate
trades and aggregation can lose individual ordering:
https://github.com/binance/binance-public-data. Validate actual file schema,
identity and data quality before coding to it.

## Size, terms and candidate split

HEAD lengths for the eleven monthly files, including partial May and without
warm-up extras: quotes 31.579 GB development and 21.936 GB holdout; individual
trades 6.478 GB development and 3.321 GB holdout. Total **63.314 GB decimal
compressed**, plus roughly 0.016 GB for all checked klines and mark files and
negligible funding. Aggregate trades would be additional, not a replacement
under the current method. Uncompressed quote sample expands from 53.8 MB to
425.8 MB (~7.9× for that partial day); extrapolating this single sample is not a
reliable total. A streaming decoder reduces working space, but plan hundreds of
GB if retaining expanded quotes and allow for source copies, checksums and backup.
Network, disk and processing capacity have not been benchmarked.

A coverage-motivated candidate is development `[2023-06-01, 2024-01-01)` and
holdout `[2024-01-01, 2024-04-01)` UTC (91 days, leap year), with at least 500
five-minute bars before each boundary. Development warm-up would start
2023-05-30 06:20 UTC; the necessary raw file continuity is **not verified**.
This starts later than the proposed May boundary because May 16 demonstrably
starts late; it does not establish all subsequent days are complete. Avoid
selecting boundaries from performance. Holdout was not acquired or evaluated.
The existing 2027 prospective holdout, pre-start freeze, 365-day adequacy floor,
and bootstrap's hard-coded 365-day block draw cannot be silently reused for the
91-day historical interval. A past, study-untouched holdout could be selected
before its acquisition and outcome inspection, but is **not prospective in
calendar time**. The choice, warm-up sharing, inference, sample floors and
disclosure require a new explicit reviewed protocol, then freeze before holdout
acquisition; no such amendment is adopted here.

Binance's Vision Dataset Terms dated 2026-08-26 say datasets are generally
CC BY-NC-SA 4.0, subject to attribution and share-alike for redistributed
derivatives; they distinguish noncommercial research from commercial use and
require a separate written license for commercial utilization:
https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md.
Assess the intended Backtest Discipline article and distribution against the
actual terms before publication or committing data-derived materials. Public
access is not unrestricted reuse; do not assume a free commercial license.

## Execution fidelity and statistics

Confirmed signals can be computed after the second right bar closes using known
trades; sequenced quotes can support approximate decision-to-fill delay, side-aware
trigger ordering and observed-price proxies when coverage and clocks are validated.
The data do not prove quote delivery completeness, executable size, queue
priority, latency, depth or market impact. Many same-millisecond updates mean a
stable within-quote-stream order must be verified; no common cross-stream order
has been established. A conservative rule for timestamp ties and gap intervals
must be frozen, not inferred from a checksum. The proposed first quote after a
latency is a simplified observation, not an executable guarantee. If its side
has insufficient quantity, a 1 BTC fill at that price is unsupported; no rule in
the current protocol resolves this. No order-book engine is recommended; instead
validate top-of-book sizes in development, select an independently justified small
size or explicit nonfill/impact policy and disclose residual limits before any
amendment. No historical fee schedule or 2023–24 instrument filters were verified;
0.05% per side is a scenario, not a demonstrated fee. Live API documentation and
current exchange info do not establish historical rules. Funding must use the
correct signed rate, settlement mark and event ordering; 1m candles cannot alone
establish the mark. Quality failure must not be masked by carry-forward prices.

Ninety-one days contain only about 13 seven-day blocks. Even 200 completed
trades would require 2.20 per day on average and would not constitute 200
independent observations. Trade rate, dependence, regime coverage and confidence
interval width are unknown; no holdout signals were inspected. A circular
block bootstrap of daily sums/counts is a reasonable candidate for clustering,
but block length, edge effects, short sample and statistical power need
independent justification and synthetic or development-only assessment. The
current 365-day floor makes the 91-day candidate automatically inadequate.
Do not lower a threshold after seeing holdout outcomes. A wide valid interval
may remain inconclusive even if mechanical data checks pass.

Alternatives: (A) the short quote-supported history may be testable only after
source, liquidity, permissions and revised inference gates pass; (B) candle-only
fills need an explicit separate protocol amendment, conservative within-bar
ordering and disclosure that bid/ask and true path are unknown; they cannot be
substituted for the current primary model; (C) stop the candidate if the owner
judges data/precision/license costs unjustified. No commercial source was acquired.

**Recommendation:** stop at operator review; investigate a bounded, checksum-
verified development-period quote/settlement sample and data-use permission, then
decide whether to authorize an explicitly amended 91-day quote-supported design.
Do not inspect or acquire holdout files, implement the full study, or freeze.

Source attribution: Binance Vision, public market-data archives and accompanying
checksums; independently computed sample diagnostics above. Dataset terms and
publication permissions must be reviewed before reuse beyond this research gate.

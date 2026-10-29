# Development acquisition and diagnostic — exploratory

The operator authorized noncommercial internal research and data acquisition;
not holdout access or freeze. Binance Vision Dataset Terms v1.0, last updated
2026-08-26, are the operator's intended noncommercial basis:
https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md.
The terms document acquired 2026-10-29 has SHA-256
`99c4921dbe023cd68c1a447775707d621453cc78aec6ac0bfe2e6c06a5114ee0`.
They describe CC BY-NC-SA 4.0 subject to their applicable conditions: attribution,
noncommercial use and share-alike for covered derivatives. Raw archive data are
not committed or redistributed. The operator says Backtest Discipline is now
noncommercial and has no ads, affiliate marketing, subscriptions or donations;
this is an operator statement, not an independent finding. Publication requires
an artifact-specific attribution and license review, including whether each
chart, table or result is a covered derivative, and applicable REST terms. No
individualized Binance approval or endorsement is claimed. Public REST access
has not been shown to grant distinct publication rights; defer that check to
publication, not to internal noncommercial development.

The tracked `provenance/development.json` records exact public locators, UTC
acquisition timestamps, sizes and SHA-256 for 41 candle archives, their checksum
sidecars, two public funding API response pages and the terms. The files reside
in ignored `data/development/` and were acquired between 08:45:20 and 08:46:10
UTC on 2026-10-29. Aggregate archive/API/terms payload size is 5,392,051 bytes,
excluding sidecars. Publisher checksum matched each candle archive; REST pages
have local SHA-256 but no publisher checksum. A full December 2023 monthly file
supplied the development warm-up; 2024 Jan–Nov monthly and Dec 1–29 daily files
supplied development. The 2024 Dec 30 file was not downloaded because it would
contain reserved holdout warm-up. Development ends 2024-12-30 00:00 UTC; the
unused 6h20m before the reserved warm-up is a file-isolation choice made before
examining development returns. Holdout `[2025-01-01,2026-01-01)` and the 500
preceding warm-up bars remain unacquired and unexamined.

Preliminary validation plus checked-in `cmd/development` confirmed 105,332
contiguous five-minute records from 2023-12-30 06:20 UTC to 2024-12-29 23:55
UTC: 500 warm-up and 104,832 development bars; no duplicates, gaps, bad OHLC,
nonfinite prices or nonconforming close timestamps were detected. It checked
archive/member names, source hashes and publisher checksums. Funding pages held
1,092 distinct records with finite rates and positive marks, one per expected
8h period; event offsets after nominal time were 0–15 ms. Validation of these
bytes and continuity does not establish exchange data completeness or live
executability. A fresh independent semantic verification is still required.

Development-only deterministic diagnostic (not a scientific claim): 1,959
trades, 719 wins, gross sum 37.66391397 R, net sum -1130.89662412 R,
mean net -0.57728261 R, net profit factor 0.40514585, four stop/target
ambiguous trades and zero recorded funding-ambiguous trades; no open trade at
boundary. Simulated 1 BTC aggregate fees 128624.35953760 USDT, adverse
slippage 51449.74438721 USDT, signed funding cost 322.94113567 USDT. A
seven-day block bootstrap diagnostic interval was [-0.63676846, -0.51860417] R;
this is exploratory, not a freeze or holdout finding. Fixed-risk notional
wealth crossed zero, so percent drawdown is not meaningful. No strategy or
acceptance setting was changed in response.

A focused independent review found the exact event-time funding approach to be
an internally coherent conservative proxy, not verified account debits, and
identified need for strict source and member identity, bootstrap tests, reporting,
tracked provenance and labeling exit times as candle identifiers. These were
addressed where stated; inferential calibration, complete reporting and semantic
replay remain unfinished. The bootstrap has constant/empty and nonconstant
independent-reference tests, but block-length suitability and coverage under
serial dependence have not been validated. Do not freeze on these diagnostics.

# Source reconstruction and feasibility — exploratory, 2026-10-28

## Evidence and identity

The user supplied playable `/home/mike/Downloads/moneyzg-scalping.mp4` (47,188,875
bytes, SHA-256 `681ed1b227a232db5af80477a6e572691abb5acf54a140574d2dcf79fe9a902a`)
and `/home/mike/Downloads/moneyzg-scalping.en.vtt` (212,206 bytes, SHA-256
`2ff56c52d899891ef2fa9cf211a1933262a0228aac3cdba842f2ab4701b45f1d`). Both
remain outside Git in Downloads; they were supplied locally, not downloaded by
this study. Their original acquisition method, publisher identity and licensing
are not independently documented. The MP4 is 19:24.981, H.264 1920×960 with
Opus audio; its container metadata does not establish a YouTube ID. The VTT is
English caption text with word timestamps; its authorship and accuracy are not
independently established. Visual frames and captions support the following
reconstruction; avoid treating captions alone as infallible. Hashes establish
local byte identity, not publisher authenticity. Check retention permissions and
backups before relying on these files for a future lineage.

The content, chapter order, MoneyZG presenter's introduction and subject are
consistent with the indexed description for
https://www.youtube.com/watch?v=nFYE0H-Lfzg, titled “Easy 5 Minute Scalping
Strategy (Simple to Follow Scalping Trading Strategy)”, MoneyZG, 2021-05-14.
The indexed snapshot reported 335,893 views. This is strong content corroboration,
not independent proof that the MP4 is an unaltered publisher download. Original
page bytes and download provenance remain unavailable. The earlier report's
inaccessible-video finding is superseded for the supplied local copy; online
playback itself was not established.

## Source-derived rules and choices

Times below are video offsets, approximate to seconds. These describe the source,
not a completed executable specification.

- **Chart and markets (00:05–01:09, 01:48–02:11):** five-minute chart; creator
  says crypto, forex or stocks and suggests other intervals are possible. He uses
  futures to demonstrate long and short and says spot can be long-only. The
  visible Binance TradingView chart says `BTCUSDT, 5` (03:25); exact contract,
  price type and historical candle provenance are not established.
- **EMAs (02:31–03:40; 03:25 screen):** 20-period yellow and 50-period red,
  using close, with chart legend `EMA (20, close, 0)` and `EMA (50, close, 0)`;
  the setting screen shows source close and offset zero for the 20. Creator
  corrects “day” to period. EMA seed, warm-up, exact implementation/version and
  numerical equality behavior are not specified by the video.
- **Fractal (03:51–04:13, 06:27–07:59):** local extreme of surrounding candles;
  creator states two more candles must complete before the entry trigger. Settings
  show `Periods 2` (06:57), unchanged default. He changes styling: a low/down
  fractal becomes a green up triangle below the bar for a potential long; a
  high/up fractal becomes a red down triangle for a potential short. The video
  does not show indicator source code or tie policy. The five-candle illustration
  at 04:48–06:10 and stated two-candle delay support the conventional two-sided
  interpretation, but do not establish every edge case or live update behavior.
- **Trend and trigger (05:11–06:17, 08:08–09:24):** long only with 20 above
  50, green low-fractal between the EMAs (below 20 and above 50); short with 20
  below 50, red high-fractal between them (above 20 and below 50). The creator
  speaks of the indicator on the candle, and chart graphics place triangles
  beyond candle extremes. Whether “between” means the extreme, drawn marker, or
  another candle price, and at which bar's EMA values, is not defined exactly.
  EMA equality, crossings during confirmation, touches and signal expiry are open.
- **Entry and exits (08:20–09:24, 10:19–10:54):** wait for signal; stop at
  50-period EMA; target one to two times stop distance, with market-dependent
  discretion. He uses approximately 1.5 in examples and mentions 1.25–1.5
  (11:54–12:10), while suggesting up to 2 for crypto. Risk/reward means price
  distance as illustrated. Whether the EMA stop is fixed at entry or follows the
  EMA is not explicitly resolved, though drawings look fixed. Entry order type,
  actual fill and latency are not specified; at 17:45–18:10 he discusses quick
  order placement, including limit and market orders, without selecting one.
- **Risk and overlap (14:32–15:33):** adding to an open winner on subsequent
  signals is optional and personal preference; he warns of increased risk. No
  fixed sizing, leverage, capital-at-risk, simultaneous position limit, netting
  or pyramiding rule is prescribed. Drawn chart-tool amounts are not a sizing
  rule. No precise session window is required; he discusses trading around an
  hour a day (18:48–19:14) as a practical option, not a defined session filter.
- **Claims and costs:** he explicitly says it does not win every time (00:27–31),
  illustrates both wins and losses (10:54–11:20, 12:32–13:38, 16:04–28), and
  claims it works better in trends (16:47–17:18). He warns fees erode small-move
  returns (18:37–48) and execution must be quick; no numerical fees, spread,
  slippage, funding or statistically measured historical win rate/return were
  found in the supplied captions or examples. Selected chart illustrations are
  not an independent performance estimate.

## Causality and illustrated entries

The creator **does not advocate knowingly trading a marker on its central candle**:
06:29–38 explicitly says two candles must complete. The source also illustrates
retrospective chart placement. For a strict five-bar pivot centred at t, the
right-hand bars are t+1 and t+2; the full predicate is available no earlier than
close(t+2), i.e. open(t)+15 minutes on continuous five-minute bars. Markers are
shown on the central extreme, not at the confirmation candle. A real trader still
needs processing and order latency after confirmation. A close and next open can
share a timestamp without proving that next-open price was executable.

- **First long (10:19–10:44; frames 10:32, 10:40, 10:45):** the green marker
  lies at the central low. He says it appears two candles afterward, then says
  enter at the start of the second candle. The position tool is visually anchored
  on the second candle to the right of the marked candle. Under the stated
  completed-two-candle rule this is one whole bar before confirmation. Frame
  geometry and the spoken instruction both support this finding; exact exchange
  timestamps and fill price are not visible.
- **Second long (11:32–11:53; frames 11:44, 11:54):** he again says wait for
  the second candle and places the entry at approximately the second candle to
  the right of the green marker. This appears to repeat the same timing error;
  chart zoom and cursor placement limit precision.
- **Losing long (12:32–13:15):** counts “one, two” before placing the position
  around the second right-hand candle (13:00–13:04). This also appears early
  under the completed-bar interpretation; it is not evidence of profitability.
- **First short (13:42–14:20; frames 14:10, 14:14):** he says enter on the
  second candle; the short position tool is anchored approximately two candles
  after the marked high. It likewise appears to precede the second right candle's
  close. Second short (15:39–16:28) is drawn retrospectively; precise marker to
  entry alignment is less clear and should not be claimed as independently proven.

Thus RQ2 finds a material mismatch between the explicitly stated confirmation
and at least the first illustrated entry. It does not prove the creator intended
an unconfirmed live strategy, nor that every illustrated fill is wrong. A
confirmation-first implementation is consistent with the verbal warning but is
not an exact reproduction of those early chart entries; disclose this as a
causal interpretation/adaptation. On an unfinished right bar the provisional
predicate may be invalidated by a new extreme. After both right bars close, a
fixed strict five-bar predicate remains stable unless input bars are revised;
actual indicator source, equality behavior and live updates remain unverified.

The checked-in synthetic tests cover strict highs and lows, t+2 availability and
provisional invalidation. They demonstrate the mathematical case only; they do
not establish TradingView's exact tie policy or test video-specific indicator
code. No noncausal result is a trading performance estimate.

## Third-party comparison

Zippy111's 2022 page
https://www.tradingview.com/script/9tAEgAaL-5min-Williams-Fractals-scalping-3commas/
says it is inspired by this video. Its reversed fractal use, default marker
between EMAs, slow-EMA stop and variable risk multiple broadly resemble the
source; its optional current-price filter, added filters and one order per trade
are implementation choices, whereas the video leaves adding positions optional.
No script source or behavior was independently verified and it is not authority
for indicator ties, exact original rules or execution timing.

## Verdict and next gate

RQ1: the strategy family is recoverable, but a unique exact executable strategy
is not: discretionary target, position handling, market selection and timing,
marker comparison, EMA seed and fills are material. RQ2: first example's entry
conflicts with the stated completed-bar confirmation; other examples appear
similar. RQ3/RQ4 remain unanswered. A prospective causal interpretation can be
scientifically meaningful if choices are independently justified and disclosed,
not advertised as an exact reproduction. Do not select choices from outcomes.

The platform supports hash-bound source manifests, post-freeze holdout binding,
immutable sealing, integrity verification and study-owned semantic replay. The
previous studies contain useful source/checksum and timing patterns, not a ready
scalping engine. A study-owned evaluator, source decoder, reproducible EMA/fractal
semantics, event/fill engine, costs and inference would still be needed. OHLCV
alone cannot establish intrabar ordering, executable bid/ask, liquidity or
latency. A shortable instrument requires funding and contract-price semantics.
Synthetic development of timing and execution primitives is technically feasible;
there is not yet a sufficiently specified candidate to implement as the primary
strategy or present for freeze review. No market or holdout data were acquired.

Next gate: decide and document a single independently justified, confirmation-first
interpretation, including target, overlap and execution choices and data
feasibility, explicitly distinguished from the retrospective examples. Then
implement and test on synthetic fixtures before prospective freeze review; the
acceptance criterion and holdout must be defined before any holdout outcomes.

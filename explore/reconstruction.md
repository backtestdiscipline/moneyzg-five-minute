# Source reconstruction and feasibility — exploratory, 2026-10-28

## Evidence and limits

Candidate: https://www.youtube.com/watch?v=nFYE0H-Lfzg. Web search's indexed
YouTube description attributes the title “Easy 5 Minute Scalping Strategy
(Simple to Follow Scalping Trading Strategy)” to MoneyZG, publication date
2021-05-14 and 335,893 views in its snapshot. This corroborates identity and
approximate candidate-selection audience, not current view count. The description
says crypto or forex, five-minute chart, EMAs, Williams fractals, and long and
short positions. A direct web open failed; playable video, original page bytes,
captions and visual evidence were not obtained. Search output is a secondary
indexed representation, not a byte-identified source; accessibility of playback
is unverified. No source bytes were acquired or hashed. No claim is made that
captions do or do not exist.

Indexed description chapters: 00:00 introduction; 00:40 scalping intro;
02:31 EMAs; 03:51 Williams fractals; 08:00 rules; 09:26, 11:29 and 12:32 long
examples; 13:40 short example; 14:32 increasing positions; 15:39 second short
example; 17:19 advice. These are navigation labels, not verified timestamps for
rules. No video-specific rule or performance claim can be timestamped until
playback and/or a reliable timestamped transcript is inspected.

| Required source item | Evidence available / open question |
| --- | --- |
| EMA periods, inputs, seed, calculations | Unknown; inspect 02:31 and chart settings. |
| Fractal definition, ties, parameters, indicator code | Unknown; inspect 03:51 and actual settings. |
| Chart and candles | Description says five minutes; exchange, instrument and price source unknown. |
| Trend conditions, long and short triggers | Description says both directions; precise conditions unknown. |
| Confirmation, entry timing, price and order | Unknown; inspect 08:00 and examples frame by frame. |
| Stop, target and exit | Unknown. |
| Position size and pyramiding | Unknown; 14:32 is only a navigation label. |
| Eligible instruments and sessions | Description mentions crypto or forex; exact instruments and restrictions unknown. |
| Original performance claims | Not established; no win rate or return inferred. |
| Fees, spread, slippage, funding | Not established; absence from the indexed description does not prove absence from video. |

No inferred detailed rule is accepted as an original rule. Thus RQ1 remains
unanswered and the examples have not been examined sufficiently to decide RQ2.

## Independent comparison, not source authority

Zippy111's 2022 TradingView page
https://www.tradingview.com/script/9tAEgAaL-5min-Williams-Fractals-scalping-3commas/
says its script was inspired by this video. Its author describes reversed fractal
signals, default marker-between-EMAs filtering, an optional current-price filter,
slow-EMA stop, a multiple-of-distance target, extra filters and one order per
trade. These are statements about that third-party implementation, not verified
original rules or a verified implementation; differences from the video cannot be
established without inspecting the video. Its screen settings are not evidence
of the video's parameters.

## Timing investigation, conditional on conventional five-bar definition

For a strict high at central bar t, high[t] exceeds highs at t-2, t-1, t+1,
t+2; reverse the comparison for a low. With finalized five-minute bars indexed
by opening time, the second right bar finishes at open(t)+15 minutes. The full
predicate is knowable only at that close; a chart may retrospectively place its
marker at t, ten minutes earlier than the open of t+2 and fifteen minutes before
confirmation. The actual displayed offset and equality policy are unknown. A
provisional intrabar indication may disappear as t+2's high/low changes; under
fixed finalized bars and a fixed strict five-bar rule the confirmed predicate
cannot later change, unless data are revised or the indicator has other logic.
Other indicator configurations can behave differently. A marker is not an order
or an executable price; next bar open has the same nominal timestamp as the
confirmation close, so execution there is an idealization needing independent
justification. The examples' causal status and whether a delayed interpretation
matches the creator's instructions remain unknown. Label any delayed version an
adaptation unless source evidence establishes equivalence. A noncausal marker
backdating exercise could illustrate timing error only, never primary trading
performance.

Synthetic tests in this directory exercise a strict five-bar illustration, both
high and low, including invalidation before close and all emitted signals' timing.
They do not verify the video's indicator, markers, entries, or future strategy.

## Platform and feasibility

Platform provides draft source manifests, hash-checked source reads, separate-data
freeze and post-freeze binding, immutable sealing, outcomes, offline integrity,
semantic-verification hooks and handoff. It does not provide a trading evaluator.
The completed hourly study contains study-owned archive/checksum validation and
OLS/HAC; the completed same-bar study demonstrates explicit availability-time
reasoning and semantic replay. Those are examples, not a generic executable
scalping engine or proof their methods suit this study. Study-owned deterministic
indicator, source decoder, event simulation, costs, inference, tests and an
independent semantic review remain necessary. A small implementation may suffice
once rules are known, but missing execution-price and intrabar path evidence may
require different data or an explicitly conservative model. Feasibility is
conditional, not proven. See the unfrozen design outline in charter.md.

Recommended next action: obtain and inspect the original video and available
captions, with permitted byte-identifiable evidence, particularly the indicator
settings and demonstrated entry times. Then decide whether the original is
reconstructible, or whether a clearly distinguished causal adaptation warrants a
new prospective question. Do not reject on the mere presence of fractals.

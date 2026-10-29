# Rights and release review

This is an evidence assessment, not legal advice. **Do not publish the empirical
summary, chart data, numerical result or article until the funding API and
archive-term applicability questions below are resolved for the exact use.**
Noncommercial operation alone does not settle them.

## Binance Vision archives

The operator reports that Backtest Discipline has no ads, affiliates, paid
subscriptions or donations. That statement has not been independently audited.
[Binance Vision Dataset Terms v1.0](https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md),
last updated 2026-08-26, list noncommercial open-source data science and
non-monetized education among allowed categories and say derived data are
covered. They specify CC BY-NC-SA 4.0, attribution to Binance Vision and the
same license for redistributed derivatives; commercial use requires a separate
written license. They prohibit hosting the datasets or derivative data feeds in
section 4.3. This project will not redistribute raw archives or API payloads,
or host a data feed; this conservative project boundary is not a claim that
every form of raw redistribution is categorically prohibited by CC itself.

The archived terms bytes acquired for the study on 2026-10-29 have SHA-256
`99c4921dbe023cd68c1a447775707d621453cc78aec6ac0bfe2e6c06a5114ee0`.
The 2025 holdout archives were acquired 2026-10-29 after that stated update.
Section 1.4 says amendments apply prospectively to datasets made available
after an amendment and preserve rights to already-downloaded datasets. It is
not established how that language applies to older 2024/2025 archives first
downloaded in 2026. Identify the effective terms and relevant Binance entity
for these exact acquisitions and territory. Do not assume first-download date
alone resolves it.

When permission and the applicable license are established, the notice should
credit Binance Vision, link the dataset and applicable terms, identify that the
study transformed the data into simulated trade and aggregate results, link
[CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/), retain
relevant copyright and disclaimer notices, and apply the same license to covered
adaptations. Do not imply Binance endorsement. Determine per artifact whether
it is an adaptation, whether other third-party rights apply, and whether its
license can actually be granted. The CC license requires attribution, a license
link, change indication, noncommercial use, share-alike where applicable and no
additional restrictions; it does not confer other necessary permissions.

## Funding REST API — separate unresolved question

The funding-history responses came from keyless `fapi.binance.com/fapi/v1/fundingRate`.
[Binance's endpoint documentation](https://developers.binance.com/docs/derivatives/usds-margined-futures/market-data/rest-api/Get-Funding-Rate-History)
documents fields and access, not a license to redistribute derived research.
Vision terms mention associated endpoints but do not clearly establish that
this REST endpoint is one. [API Keys Product Terms v1.0](https://bin.bnbstatic.com/static/cms/cg08ou2ak0tn7mcplvfg/file/826cc99cc3847cd017dc7b82525896f01ac211d9c177573d8ea958402766d965.pdf),
effective 2026-01-05, address API keys; the study used no key, and those terms
alone do not establish rights for this endpoint. Determine the applicable
entity, terms and permission for publishing aggregated outcomes calculated
using funding rates and marks. Omission of raw responses alone does not resolve
rights in a derived result.

## Release decision

Narrowest supported immediate scope: use the package privately to prepare and
review a factual description of the public video and proposed methodology;
public use of source video, captions or screenshots has its own rights review.
Publication of aggregate results, charts or explanatory text that incorporates
the outcome is conditional on resolving the two source-rights questions and
applying any required notices. No individualized permission or endorsement has
been obtained. Do not infer permission from availability, the lack of an
express prohibition, or a noncommercial label. Once cleared, publish only the
specific files marked conditional in `release-manifest.json`, or their
individually reviewed excerpts, under the terms determined to apply.

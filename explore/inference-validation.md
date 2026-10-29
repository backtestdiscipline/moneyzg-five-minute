# Statistical inference repair — exploratory, 2026-10-29

Classification: INFERENCE_VALIDATED **for the prespecified synthetic scenarios**;
this is not a guarantee of valid coverage for arbitrary market dependence or
freeze readiness. Candidate and test plan were committed as `ef1662c` before
the final simulation. The earlier seven-day failure remains in
`final-development-audit.md`; no strategy or cost parameter was changed.

## What failed and what changed

The 7-day circular-block percentile bootstrap resampled paired daily net-R
sums and trade counts, including zero days, so the unit and ratio estimand were
appropriate. Seven-day blocks truncate cross-block dependence and produce only
approximately 52 blocks in 365 days. Under persistent AR(0.8), estimated
uncertainty was too narrow, yielding about 81% coverage. Varying independent
daily trade frequency produced a similar shortfall, so variable counts alone
are not the demonstrated cause. The block-length sensitivity supports a
short-block/dependence explanation; the current tests do not separately isolate
percentile interval construction or effective sample size. Independent
reference tests found no demonstrated implementation error. Longer blocks alone
are not a universal remedy; the old simulations' 28-day percentile bootstrap
still under-covered. Serial dependence, nonstationarity, small effective sample
and interval construction all remain possible contributors in real data.

The candidate uses five chronological 73-day batches of paired daily sums and
counts, ratio theta = sum(S)/sum(N), batch influence sums
`Z_j = sum_j(S-theta*N)`, standard error `sqrt(5/4 * sum(Z_j²))/sum(N)`, and
two-sided interval `theta ± 2.7764451051977987*SE`. It retains zero-trade days,
actual counts and all magnitudes. Trades crossing days, including their costs
and funding, are assigned to UTC exit day. No seed or replications are used in
the interval itself. The t(4) critical value and five batches make this a
small-sample approximation assuming sufficiently weak dependence across batches;
five batches cannot reliably diagnose long-memory, regime shifts or infinite
variance. See `inference-preregistration.md` for the complete a-priori method,
admission, failure and scenario definitions. Holdout adequacy retains 365 full
UTC days, >=200 trades, >=30 exit days and no terminal position, and adds >=1
trade in each batch. Invalid data or variance fail technically. Positive modeled
expectancy still requires the lower 95% interval bound >0.

## Final synthetic suite

The independently coded simulations use 3,000 fixed-seed independent 365-day
series per scenario, exactly the prespecified acceptance floor of 92% in each
and zero unexpected errors. Wilson limits below are binomial 95% Monte Carlo
intervals, not the strategy interval.

| Process | Covered | Coverage | Wilson 95% | MCSE |
| --- | ---: | ---: | ---: | ---: |
| Independent Gaussian | 2850/3000 | 95.00% | 94.16–95.72% | 0.40 pp |
| Independent Student t(3) | 2848/3000 | 94.93% | 94.09–95.66% | 0.40 pp |
| Gaussian AR(1), phi .5 | 2831/3000 | 94.37% | 93.48–95.14% | 0.42 pp |
| Gaussian AR(1), phi .8 | 2836/3000 | 94.53% | 93.66–95.29% | 0.42 pp |
| AR(1) .8, counts uniform 0–4 | 2832/3000 | 94.40% | 93.52–95.17% | 0.42 pp |
| Clustered volatility | 2857/3000 | 95.23% | 94.41–95.94% | 0.39 pp |

All six pass the declared floor, with no numerical errors. The 92% floor is a
validation threshold, not a new confidence level or a scientific success rule.
The tests do not show nominal coverage for untested data-generating processes.
Code: `explore/inference-batch/main.go`; run offline with Go 1.27.2. An
independent Python implementation `explore/inference-reference/reference.py`
uses a different fixed PRNG and independently reproduced ~94–96% coverage for
all corresponding cases; its volatility scenario is an additional variant,
not the exact preregistered final suite. Its separate variable-count reference
fixture agrees with the Go result to 1e-14. The Go tests also check a closed-form
fixture, malformed input, zero variance, sample gates and deterministic output.

## Development-only application

The development partition has 364 full days, whereas the primary method fixes
365. Without inventing a day or changing holdout rules, the diagnostic uses
the first 360 chronological days in five 72-day batches; the last four days
are excluded solely by this calendar rule. This is a related exploratory
calculation, not the exact 365-day primary or a holdout result. Its mean net R
is -0.57966802 and interval [-0.66995997,-0.48937606]. On those same 360 days,
the previous seven-day bootstrap gives [-0.63841223,-0.51996622]. The full
364-day development interval was [-0.63676846,-0.51860417] under that earlier
method. The wider batch interval on the same 360 days reflects its larger
estimated uncertainty and t critical value; it does not prove nominal coverage
on this dataset. No outcomes were used to alter trading rules.

## Remaining work

The inference candidate may proceed to independent freeze-contract review,
conditional on its limited assumptions and a clear policy for sample adequacy,
nonstationarity, and technical failure. The draft platform records still carry
an explicit inference blocker and the old proposal; they must be updated and
reviewed before a future freeze. The study's platform run/seal and semantic
verification path remain unimplemented by instruction in this iteration.
Publication-specific attribution and licensing review remain separate. No
2025 holdout bytes were acquired or examined, and no freeze was performed.

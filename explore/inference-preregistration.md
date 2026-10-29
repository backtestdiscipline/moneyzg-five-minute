# Inference candidate and validation plan — fixed before final simulation

Registered 2026-10-29 before running the present validation suite or applying
the candidate to the development ledger. This is a provisional research choice,
not an approved frozen contract and not a trading parameter change.

Primary estimand is expected net R per completed trade, estimated by the ratio
`theta = sum(S_d)/sum(N_d)`, where for each UTC exit day S_d is the sum of net
R of trades exiting that day and N_d their count; zero-trade days contribute
(0,0). A trade crossing midnight is assigned entirely to its exit day; funding
and modeled execution costs stay inside its net R. High-count days retain their
actual counts and contribution: no clipping, winsorization, equal-day weighting,
normalization, or filtering by returns.

Candidate inference: partition exactly 365 chronological days into five fixed,
nonoverlapping contiguous 73-day blocks. Let `Z_j = sum_block(S_d - theta*N_d)`.
Standard error is `sqrt((5/4)*sum(Z_j^2))/sum(N_d)` and the two-sided interval
is `theta ± 2.7764451051977987*SE`, the Student-t 97.5th percentile with four
degrees of freedom. No randomness or bootstrap draws enter the estimate. This
batch-means ratio-influence procedure treats long batches as approximately
independent, with a small-sample t approximation; it is not exact for dependent
or nonstationary observations. Fixed 73-day blocks were chosen a priori to
retain materially longer dependence than the failed seven-day method and use
all 365 days, not by looking at development outcomes. Five blocks permit
limited effective degrees of freedom and can give very wide or misleading
intervals when dependence persists across blocks.

Adequacy requires exactly 365 complete days, >=200 completed trades, >=30 exit
days, >=1 trade in each 73-day block, all source and settlement checks, and no
terminal open position. Missing data, negative counts, nonfinite values,
nonzero daily sums with zero count, nonpositive total counts, nonfinite or
nonpositive estimated variance, and arithmetic overflow are technical failures;
inadequate trade/day/block counts are INCOMPLETE. No special censoring of outliers
or irregularly active days. Success requires adequacy and the two-sided 95%
interval's lower bound >0. The estimator assumes sufficiently stable,
finite-variance weakly dependent daily influence observations; it does not
guarantee validity under structural breaks, extremely long memory or infinite
variance. A publication must disclose limitations.

Final validation suite fixed now: 3,000 independently generated datasets per
scenario, each n=365, independent Gaussian; independent finite-variance Student
t with 3 degrees of freedom; stationary Gaussian AR(1) phi 0.5 and phi 0.8;
AR(1) phi 0.8 with daily counts independent uniform {0,1,2,3,4} and daily S =
count × shock; and stationary volatility clustering from independent Gaussian
return innovations scaled by exp(0.5*v), where v is stationary Gaussian AR(1)
phi 0.8, innovation SD 0.6 (stationary SD 1). For other scenarios N=1 and S is
the generated outcome. True mean and ratio in every case are zero. Use fixed
seed, independent test data per replicate, and an independent implementation to
cross-check formulas. Capture empirical 95% interval coverage, Wilson 95%
binomial confidence intervals, MCSE, failures and mean interval width. The
predefined validation rule requires >=92.0% empirical coverage in EACH scenario
and zero unexpected numerical failures. A passing suite is evidence for these
models only, not universal validity. This threshold is not the scientific
claim's 95% confidence level; it flags material observed undercoverage with
~0.4 percentage-point Monte Carlo standard error near 95% at 3,000 replicates.
If any scenario fails, do not tune and rerun candidates; report the failure.
Only after validation passes may the candidate be applied to the already known
2024 development ledger; never access 2025 in this invocation.

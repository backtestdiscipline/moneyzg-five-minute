# Freeze record

On 2026-10-29 the operator approved the candidate in `freeze-review.md` for
freezing. The platform froze the unchanged, previously reviewed specification
at 2026-10-29T13:20:19Z with `--exploration separate-data`. The freeze discloses
exploratory 2024 development and synthetic validation; the reserved historical
2025 holdout was not acquired or evaluated. No scientific verdict exists.

- Specification SHA-256: `14d326cf92bb92fefd5ef65fcaa6fe8e074aeea75aca3f9073df86ad1b44697c`.
- Manifest SHA-256: `0f4845638d598fa41507b605b2e9450e2134011d6e0dfd618dc68e377872a7c6`.
- Freeze record SHA-256: `10d9f672f67cb469f7e776fa9aaad13930abc8ca7657654cba02c8c3f4d5a0c5`.
- Canonical lineage SHA-256: `2ff2c786e8f62c23f0211b5671194471bd44b437166f86f4b85376ff7fa47d12`.
- Frozen clean study commit: `a638d9111acaf64d72ef8f3543072c82d0c6eb93`,
  containing the evaluator, inference, semantic replay and `go.mod` with Go
  toolchain 1.27.2 and a local research-platform dependency.

The freeze record binds the study commit and manifest/spec identities. The
platform checkout used for freeze was clean at commit
`9d0cb4f2a808ae2743c5bac90b2ffeeafaec3750`; under the existing platform
contract this is contextual provenance, not a platform identity embedded in
the freeze. BeginRun records the actual clean platform implementation identity
and study implementation identity in its execution plan. The first future run
binds the still-pending post-freeze source bytes; it does not retroactively
make those bytes known at freeze.

Before freeze and again after freeze, `research-platform check
moneyzg-five-minute` passed with evaluator and inference admitted as
STUDY_OWNED. `research-platform status .` reported the lineage frozen.
SHA-256 values were independently computed from the three records. The platform
`verify .` command could not inspect a run because `.research-platform/runs`
does not exist; this is not a run integrity PASS and no run was created. All 30
post-freeze paths were absent at the pre-freeze check and carried no acquired
identity. The earlier isolated synthetic lifecycle verification is documented
in `freeze-review.md` and is not a real holdout run.

Holdout warm-up: [2024-12-30T06:20:00Z, 2025-01-01T00:00:00Z), exactly 500
five-minute bars. Evaluation: [2025-01-01T00:00:00Z, 2026-01-01T00:00:00Z),
365 calendar days. No warm-up pivot is eligible; an open final position is
INCOMPLETE. The predeclared five 73-day batches and acceptance and technical
failure criteria are in the frozen specification. This is an untouched
historical holdout, not calendar-forward evidence.

Separate operator authorization is required before acquisition and evaluation.
Future work must acquire and verify all 30 declared post-freeze artifacts,
respect source/checksum and funding requirements, keep raw data out of Git,
use a clean study worktree, and use the established run and semantic verification
path. Noncommercial licensing and eventual publication review remain as
described in `data-decision.md` and `freeze-review.md`.

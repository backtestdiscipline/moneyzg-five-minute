# MoneyZG five-minute scalping study

This is an independent study repository using the adjacent research-platform.
See `STATUS.md`, `charter.md`, and `explore/ohlcv-protocol.md` for its current
unfrozen scope and blockers. No holdout has been evaluated.

The study was initially developed under `research-platform/studies/moneyzg-five-minute`.
Its earlier study-only history was extracted into this repository with `git
subtree split`; commit IDs therefore differ from the original platform commits.
Original platform commits remain in platform Git history, ending at `9177c73`.
Source bytes remain outside Git; no sealed evidence or source bytes were moved.
The `ohlcv` package is exploratory study-owned code, not a frozen evaluator.

Run synthetic checks with Go 1.27.2 and offline module resolution:

```sh
GOPROXY=off GOSUMDB=off go test ./...
```

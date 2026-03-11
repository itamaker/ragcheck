# ragcheck

`ragcheck` gives you fast offline retrieval metrics without requiring Python notebooks.

## Example

```bash
go run . score -qrels examples/qrels.json -run examples/run.json -k 3
```

## Metrics

- `Precision@k`
- `Recall@k`
- `HitRate@k`
- `MRR@k`

This makes it suitable for both engineering smoke checks and lightweight research baselines.

## Install

From source:

```bash
go install github.com/YOUR_GITHUB_USER/ragcheck@latest
```

From Homebrew after you publish a tap formula:

```bash
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/ragcheck
```

## Repo-Ready Files

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `.goreleaser.yaml`
- `PUBLISHING.md`
- `scripts/render-homebrew-formula.sh`

## Release

```bash
git tag v0.1.0
git push origin v0.1.0
```

The tagged release workflow publishes multi-platform binaries and `checksums.txt`, which you can feed into the Homebrew formula renderer.
The generated formula should be committed to `https://github.com/itamaker/homebrew-tap`.

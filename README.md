# ragcheck

[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

`ragcheck` is a Go CLI for offline retrieval evaluation.

It gives you fast `Precision@k`, `Recall@k`, `HitRate@k`, and `MRR@k` scores without spinning up notebooks or Python-based eval tooling.

## Support

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/amaker)

## Quickstart

### Install

Install with your preferred method:

```bash
# From the custom tap
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/ragcheck
```

```bash
# Or install from source
go install github.com/itamaker/ragcheck@latest
```

<details>
<summary>You can also download binaries from <a href="https://github.com/itamaker/ragcheck/releases">GitHub Releases</a>.</summary>

Current release archives:

- macOS (Apple Silicon/arm64): `ragcheck_0.1.0_darwin_arm64.tar.gz`
- macOS (Intel/x86_64): `ragcheck_0.1.0_darwin_amd64.tar.gz`
- Linux (arm64): `ragcheck_0.1.0_linux_arm64.tar.gz`
- Linux (x86_64): `ragcheck_0.1.0_linux_amd64.tar.gz`

Each archive contains a single executable: `ragcheck`.

</details>

If the repository is still private, release-based installs require GitHub access to the repository assets.

### First Run

Run:

```bash
ragcheck score -qrels examples/qrels.json -run examples/run.json -k 3
```

## Requirements

- Go `1.22+`

## Run

```bash
go run . score -qrels examples/qrels.json -run examples/run.json -k 3
```

Supported metrics:

- `Precision@k`
- `Recall@k`
- `HitRate@k`
- `MRR@k`

## Build From Source

```bash
make build
```

```bash
go build -o dist/ragcheck .
```

## What It Does

1. Loads qrels and retrieval run files from JSON.
2. Matches retrieved document IDs against relevant sets.
3. Computes standard top-k retrieval metrics.
4. Prints quick offline evaluation summaries for engineering and research loops.

## Notes

- Keep qrels and run files aligned on `query_id`.
- Maintainer release steps live in `PUBLISHING.md`.

## Contributors ✨

| [![Zhaoyang Jia][avatar-zhaoyang]][author-zhaoyang] |
| --- |
| [Zhaoyang Jia][author-zhaoyang] |



[author-zhaoyang]: https://github.com/itamaker
[avatar-zhaoyang]: https://images.weserv.nl/?url=https://github.com/itamaker.png&h=120&w=120&fit=cover&mask=circle&maxage=7d

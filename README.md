[![License](https://img.shields.io/github/license/aaronflorey/ubuntu-mirrorselect?style=flat)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/aaronflorey/ubuntu-mirrorselect/ci.yaml?branch=main&style=flat&label=ci)](https://github.com/aaronflorey/ubuntu-mirrorselect/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/aaronflorey/ubuntu-mirrorselect?display_name=tag&sort=semver&style=flat)](https://github.com/aaronflorey/ubuntu-mirrorselect/releases)

# MirrorSelect

MirrorSelect is a command-line tool that ranks Ubuntu archive mirrors by TCP latency and download speed so you can choose a faster package mirror.

This repository is a maintained fork of the original upstream project at [haukened/mirrorselect](https://github.com/haukened/mirrorselect).

## Why MirrorSelect

- Uses the Launchpad mirror list instead of the smaller mirror redirect service.
- Measures TCP responsiveness instead of ICMP latency alone.
- Re-checks the best candidates by download speed before printing results.

## Installation

Homebrew cask (after the first tagged release):

```sh
brew install --cask aaronflorey/tap/mirrorselect
```

Source build:

```sh
git clone https://github.com/aaronflorey/ubuntu-mirrorselect.git
cd ubuntu-mirrorselect
go build -o mirrorselect .
```

GitHub releases:

- Download a prebuilt archive from the [Releases](https://github.com/aaronflorey/ubuntu-mirrorselect/releases) page.
- Run `curl -fsSL https://raw.githubusercontent.com/aaronflorey/ubuntu-mirrorselect/main/install.sh | bash` to install the latest release for the current host into `~/.local/bin`.
- If you prefer the Amber source, run `amber run scripts/install-latest-release.ab`.

## Setup

- MirrorSelect is intended for Ubuntu hosts and verifies that the running system is Ubuntu.
- On Ubuntu systems, MirrorSelect auto-detects the release codename from `/etc/os-release`, `/etc/lsb-release`, or `lsb_release -cs`.
- If you do not want GeoIP-based country detection, pass `--country` explicitly.

## Usage

```sh
mirrorselect [flags]
```

Important flags:

- `--arch`: mirror architecture to target.
- `--country`: ISO 3166-1 alpha-2 country code such as `US` or `DE`.
- `--max`: maximum number of mirrors to benchmark by download speed.
- `--interactive` / `-i`: show ranked mirrors with speed/latency and prompt to pick one.
- `--apply`: write the selected mirror into the active APT source file, preferring `/etc/apt/sources.list.d/ubuntu.sources` and falling back to `/etc/apt/sources.list`, after creating a timestamped backup.
- `--output`: `text` (default) or `json`.
- `--protocol`: `http`, `https`, or `any`.
- `--release`: Ubuntu codename such as `noble`, `jammy`, or `focal`.
- `--timeout`: latency timeout in milliseconds.
- `--verbosity`: `DEBUG`, `INFO`, `WARN`, or `ERROR`.
- `--yes`: skip confirmation prompt when using `--apply`.

Note:

- MirrorSelect only updates APT files when you pass `--apply`.
- If `--apply` is used without root privileges, MirrorSelect re-runs itself with `sudo` and forwards prompts/output.
- `--apply` writes entries for `release`, `release-updates`, `release-backports`, and `release-security`, and includes both `deb` and `deb-src` source types.

Examples:

```sh
mirrorselect
mirrorselect --protocol https --country US
mirrorselect --release noble --country GB --max 5 --timeout 800
mirrorselect --country US --protocol https --max 5 --interactive
mirrorselect --country US --output json
sudo mirrorselect --country US --interactive --apply
```

## Development

```sh
lefthook install
gofmt -w .
go vet ./...
go test ./...
```

Lefthook keeps `install.sh` compiled from `scripts/install-latest-release.ab` on `pre-commit`, and runs `go vet`, `go test`, and `go build` on `pre-push`.

## Releases

- CI runs on pushes and pull requests to `main` and `master`.
- `release-please` opens or updates the release PR from conventional commits.
- Merging the release PR creates a `vX.X.X` tag, publishes GitHub release archives, and updates the Homebrew tap cask.

## Security

Please see [SECURITY.md](SECURITY.md).

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).

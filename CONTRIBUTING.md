# Contributing to MirrorSelect

Thanks for contributing. By participating in this project, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development Setup

```sh
git clone https://github.com/aaronflorey/ubuntu-mirrorselect.git
cd ubuntu-mirrorselect
lefthook install
go test ./...
```

## Before Opening a Pull Request

Run the local checks:

```sh
gofmt -w .
amber build scripts/install-latest-release.ab install.sh
go vet ./...
go test ./...
```

Update nearby tests and documentation when behavior changes.

Lefthook uses `lefthook.yml` to rebuild `install.sh` from `scripts/install-latest-release.ab` on `pre-commit` and to run the main Go checks on `pre-push`.

## Commit Style

This repository uses conventional commits so release automation can generate release notes and version bumps.

Examples:

- `feat: add country validation`
- `fix: handle empty mirror results`
- `docs: clarify Linux-only auto-detection`

## Reporting Issues

- Bugs and feature requests: open an issue at [github.com/aaronflorey/ubuntu-mirrorselect/issues](https://github.com/aaronflorey/ubuntu-mirrorselect/issues).
- Security issues: follow [SECURITY.md](SECURITY.md) instead of opening a public issue.

## Pull Requests

1. Create a branch for your change.
2. Keep the change focused.
3. Add or update tests when needed.
4. Open a pull request with a clear summary of the problem and fix.

## License

By contributing to MirrorSelect, you agree that your contributions will be licensed under the [MIT License](LICENSE).

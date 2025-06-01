# Introduction

[![Tests](https://github.com/anttiharju/relcheck/actions/workflows/tests.yml/badge.svg)](https://github.com/anttiharju/relcheck/actions/workflows/tests.yml) [![Linters](https://github.com/anttiharju/relcheck/actions/workflows/linters.yml/badge.svg)](https://github.com/anttiharju/relcheck/actions/workflows/linters.yml) [![Docs build](https://github.com/anttiharju/relcheck/actions/workflows/docs-build.yml/badge.svg)](https://github.com/anttiharju/relcheck/actions/workflows/docs-build.yml)

## Why

1. Documentation is useful; documentation with broken relative links is less so.
2. `mkdocs build --strict` is too strict, can not check files outside of `docs/` directory.
3. Other existing tools I found were too slow, taking up to 10 seconds. This tool typically runs in milliseconds:

```sh
$ hyperfine "relcheck run"
Benchmark 1: relcheck run
  Time (mean ± σ):      32.2 ms ±   0.2 ms    [User: 11.7 ms, System: 15.2 ms]
  Range (min … max):    31.3 ms …  32.8 ms    84 runs
```

## Installation

```sh
sudo sh -c "curl -sSfL https://raw.githubusercontent.com/anttiharju/relcheck/HEAD/relcheck.bash -o /usr/local/bin/relcheck && chmod +x /usr/local/bin/relcheck"
```

Note: the tool depends on awk, and not all versions of awk are apparently compatible. Install `gawk` in case you're having issues.

Eventually there will be a rewrite to produce a static binary without this issue.

## Usage

In integrated terminals of editors such as VS Code, the reported broken links such as `dist/brew/README.md:5:19` are clickable when holding ctrl/cmd to bring your cursor right to where the ^ indicator points:

```sh
$ relcheck run
dist/brew/README.md:5:19: broken relative link (file not found):
- [`values.bash`](./values.sh) is required by the [render-template](https://github.com/anttiharju/actions/tree/v0/render-template) action.
                  ^
```

The `file:line:column` link syntax is the same one that golangci-lint uses.

### Manual

Using defaults inside a Git repository

```sh
relcheck run
```

for advanced usage, refer to the printed out info from

```sh
relcheck
```

### Git pre-commit hook (via Lefthook)

[Lefthook](https://github.com/evilmartians/lefthook) is an awesome Git hooks manager, enabling [shift-left testing](https://en.wikipedia.org/wiki/Shift-left_testing) that improves developer experience. `relcheck` was built for usage with Lefthook. Here is a minimal `lefthook.yml` example:

```yml
output:
  - success
  - failure

pre-commit:
  parallel: true
  jobs:
    # Install from https://github.com/anttiharju/relcheck
    - name: relcheck
      run: relcheck run
```

### GitHub Actions

A composite action is available through my [actions monorepo](https://github.com/anttiharju/actions/tree/v0/relcheck). Here is a minimal `.github/workflows/build.yml` example:

```yml
name: Build
on:
  push:
    branches:
      - main
  pull_request:

jobs:
  validate:
    name: Validate
    runs-on: ubuntu-24.04
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: relcheck
        uses: anttiharju/actions/relcheck@fa0a8b6cd47e30e4abf7ce4fbbd8ec0f377405db
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/anttiharju/relcheck.svg?variant=adaptive)](https://starchart.cc/anttiharju/relcheck)

## Why

This tool was developed alongside and mainly for https://github.com/anttiharju/vmatch. Idea is to have `vmatch` to be as linted as possible to make maintaining the project a breeze. Additionally the tooling built to support it will make my future projects easier to work on, allowing me to mostly focus on functionality without existing things breaking while I refactor the projects to my will.

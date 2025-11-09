# Introduction

[![Build](https://github.com/anttiharju/relcheck/actions/workflows/build.yml/badge.svg?event=push)](https://github.com/anttiharju/relcheck/actions/workflows/build.yml)

## Why

1. Documentation is useful; documentation with broken relative links is less so.
2. `mkdocs build --strict` is too strict, can not check files outside of `docs/` directory.
3. Other existing tools I found were too slow, taking up to 10 seconds. This tool typically runs in milliseconds:

```sh
$ hyperfine "relcheck all"
Benchmark 1: relcheck all
  Time (mean ± σ):       7.2 ms ±   0.5 ms    [User: 3.8 ms, System: 3.2 ms]
  Range (min … max):     6.5 ms …   9.7 ms    286 runs
```

## Installation

```sh
brew install anttiharju/tap/relcheck
```

Or download a binary from a GitHub release.

### Updating

```sh
brew update && brew upgrade relcheck
```

## Usage

Using defaults inside a Git repository

```sh
relcheck all
```

for advanced usage, refer to the printed out info from

```sh
relcheck
```

Although the recommendation is to setup a integration via Lefthook or GitHub Actions instead of manual use.

### GitHub Actions

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
        uses: anttiharju/relcheck@v1.8.11
```

## Integrations

### VS Code

The reported broken links such as `dist/brew/README.md:5:19` are clickable in the intergrated terminal when holding ctrl/cmd. It will bring you right to where the ^ indicator points:

```sh
$ relcheck all
dist/brew/README.md:5:19: broken relative link (file not found):
- [`values.bash`](./values.sh) is required by the [render-template](https://github.com/anttiharju/actions/tree/v1/render-template) action.
                  ^
```

The `file:line:column` link syntax is the same one that golangci-lint uses.

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
      run: relcheck all
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/anttiharju/relcheck.svg?variant=adaptive)](https://starchart.cc/anttiharju/relcheck)

## Why

This tool was developed alongside and mainly for https://github.com/anttiharju/vmatch. Idea is to have `vmatch` to be as linted as possible to make maintaining the project a breeze. Additionally the tooling built to support it will make my future projects easier to work on, allowing me to mostly focus on functionality without existing things breaking while I refactor the projects to my will.

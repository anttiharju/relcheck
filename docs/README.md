# Introduction

[![Tests](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/tests.yml/badge.svg)](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/tests.yml) [![Linters](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/linters.yml/badge.svg)](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/linters.yml) [![Docs build](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/docs-build.yml/badge.svg)](https://github.com/anttiharju/check-relative-markdown-links/actions/workflows/docs-build.yml)

## Why

1. Documentation is useful; documentation with broken relative links is less so.
2. `mkdocs build --strict` is too strict, can not check files outside of `docs/` directory.
3. Other existing tools I found were too slow, taking up to 10 seconds. This tool typically runs in milliseconds:

```sh
$ time check-relative-markdown-links run
________________________________________________________
Executed in   32.02 millis    fish           external
  usr time   11.29 millis    0.23 millis   11.05 millis
  sys time   16.60 millis    1.87 millis   14.73 millis
```

## Installation

```sh
sudo sh -c "curl -sSfL https://raw.githubusercontent.com/anttiharju/check-relative-markdown-links/HEAD/check-relative-markdown-links.bash -o /usr/local/bin/check-relative-markdown-links && chmod +x /usr/local/bin/check-relative-markdown-links"
```

## Usage

In integrated terminals of editors such as VS Code, the reported broken links such as `dist/brew/README.md:5:19` are clickable when holding ctrl/cmd to bring your cursor right to where the ^ indicator points:

```sh
$ check-relative-markdown-links run
dist/brew/README.md:5:19: broken relative link (file not found):
- [`values.bash`](./values.sh) is required by the [render-template](https://github.com/anttiharju/actions/tree/v0/render-template) action.
                  ^
```

The `file:line:column` link syntax is the same one that golangci-lint uses.

### Manual

Using defaults inside a Git repository

```sh
check-relative-markdown-links run
```

for advanced usage, refer to the printed out info from

```sh
check-relative-markdown-links
```

### Git pre-commit hook (via Lefthook)

[Lefthook](https://github.com/evilmartians/lefthook) is an awesome Git hooks manager, enabling [shift-left testing](https://en.wikipedia.org/wiki/Shift-left_testing) that improves developer experience. `check-relative-markdown-links` was built for usage with Lefthook. Here is a minimal `lefthook.yml` example:

```yml
output:
  - success
  - failure

pre-commit:
  parallel: true
  jobs:
    # Install from https://github.com/anttiharju/check-relative-markdown-links
    - name: check-relative-markdown-links
      run: check-relative-markdown-links run
```

### GitHub Actions

A composite action is available through my [actions monorepo](https://github.com/anttiharju/actions/tree/v0/check-relative-markdown-links). Here is a minimal `.github/workflows/build.yml` example:

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

      - name: check-relative-markdown-links
        uses: anttiharju/actions/check-relative-markdown-links@fa0a8b6cd47e30e4abf7ce4fbbd8ec0f377405db
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/anttiharju/check-relative-markdown-links.svg?variant=adaptive)](https://starchart.cc/anttiharju/check-relative-markdown-links)

## Why

This tool was developed alongside and mainly for https://github.com/anttiharju/vmatch. Idea is to have `vmatch` to be as linted as possible to make maintaining the project a breeze. Additionally the tooling built to support it will make my future projects easier to work on, allowing me to mostly focus on functionality without existing things breaking while I refactor the projects to my will.

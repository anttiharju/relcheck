# check-relative-markdown-links

Documentation is useful; documentation with broken links to other files in the repository is less so.

## Installation

Run

```sh
curl -sSfL https://raw.githubusercontent.com/anttiharju/check-relative-markdown-links/HEAD/install.sh | sh
```

and

```sh
sudo ln -sf ~/.anttiharju/check-relative-markdown-links.bash /usr/local/bin/check-relative-markdown-links
```

## Usage

Using defaults inside a Git repository

```sh
check-relative-markdown-links run
```

for advanced usage, refer to the printed out info from

```sh
check-relative-markdown-links
```

### Lefthook

[Lefthook](https://github.com/evilmartians/lefthook) is an awesome Git hooks manager. It enables [shift-left testing](https://en.wikipedia.org/wiki/Shift-left_testing), improving developer experience. `check-relative-markdown-links` was built for usage with Lefthook. Here's a minimal example `lefthook.yml` configuration file:

```yml
output:
  - success
  - failure

pre-commit:
  parallel: true
  jobs:
    - name: check-relative-links
      run: check-relative-links run
```

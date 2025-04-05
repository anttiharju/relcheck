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

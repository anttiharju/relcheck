# Issues caught

## Preface

This tool ignores .markdown files by default, so you can still retain syntax highlighting in your editor and demonstrate issues caught by the tool in real markdown.

If you want to use file extensions other than `.md`, you can feed any desired files to the tool with a command such as

```sh
check-relative-markdown-links --verbose $(git ls-files '*.markdown')
```

although be aware that ShellCheck will complain of such use. There's a way to do it in a compliant way, but it won't be a neat oneliner. Also files with spaces won't work with the above approach, as ShellCheck correctly points out with [SC2046](https://www.shellcheck.net/wiki/SC2046). ShellCheck is awesome.

## Links

Broken links, such as typos are caught [../REDME.md](../REDME.md).

## Anchors

Similarly non-existent anchors are also caught [README.md#gitlab-actions](../README.md#gitlab-actions)

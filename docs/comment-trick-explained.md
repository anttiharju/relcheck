# Comment trick

TL;DR: update the url with the change you did to the commented out relative link.

These two file paths should match:

```
<!--[README](../testdata/happy/README.md)-->
[testdata/happy/README.md](https://github.com/anttiharju/check-relative-markdown-links/blob/HEAD/testdata/happy/README.md)
```

## Problem

If you build your documentation site with [MkDocs](https://www.mkdocs.org) (which, btw, if you use [Backstage](https://backstage.io), you do) you may have found out that making relative links out of the `docs/` directory do not work on the final site. `mkdocs build --strict` displays a `WARNING` about this:

```sh
mkdocs build --strict
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /Users/antti/anttiharju/check-relative-markdown-links/site
WARNING -  Doc file 'comment-trick.md' contains a link '../testdata/happy/README.md', but the target is not found among documentation files.

Aborted with 1 warnings in strict mode!
```

<!--[README](../testdata/happy/README.md)-->

So as a workaround you can link to your GitHub-hosted Markdown file like this: [testdata/happy/README.md](https://github.com/anttiharju/check-relative-markdown-links/blob/HEAD/testdata/happy/README.md) and `mkdocs build --strict` is happy again, yay!

```
mkdocs build --strict
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /Users/antti/anttiharju/check-relative-markdown-links/site
INFO    -  Documentation built in 0.14 seconds
```

But with this workaround you have lost static checking that makes refactors annoying because you'll get hidden breakage.

## Solution (workaround)

Add the relative link, but comment it out. This way you still get a tripwire for refactors but MkDocs is happy.

```
<!--[README](../testdata/happy/README.md)-->
```

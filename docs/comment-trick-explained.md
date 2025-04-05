# Comment trick

## TL;DR

For a reason a relative link has to be an URL instead and someone decided to introduce static checking by burying the relative link into a comment somewhere nearby the URL actually used in the docs, for example like this:

[../testdata/README.md](../testdata/README.md)

```md
<!--[../testdata/README.md](../testdata/README.md) https://anttiharju.dev/check-relative-markdown-links/comment-trick-explained-->

MkDocs disallows relative links outside of the docs directory, so here's a GitHub one instead: https://github.com/anttiharju/check-relative-markdown-links/blob/HEAD/testdata/README.md
```

## Problem

If you build your documentation site with [MkDocs](https://www.mkdocs.org) (which, btw, if you use [Backstage](https://backstage.io), you do) you may have found out that making relative links out of the `docs/` directory do not work on the final site. `mkdocs build --strict` displays a `WARNING` about this:

```sh
$ mkdocs build --strict
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /Users/antti/anttiharju/check-relative-markdown-links/site
WARNING -  Doc file 'comment-trick-explained.md' contains a link '../testdata/README.md', but the target is not found among documentation files.

Aborted with 1 warnings in strict mode!
```

<!--[testdata/README.md](../testdata/README.md) https://anttiharju.dev/check-relative-markdown-links/comment-trick-explained-->

So as a workaround you can link to your GitHub-hosted Markdown file like this: [testdata/README.md](https://github.com/anttiharju/check-relative-markdown-links/blob/HEAD/testdata/README.md) and `mkdocs build --strict` is happy again, yay!

```sh
$ mkdocs build --strict
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /Users/antti/anttiharju/check-relative-markdown-links/site
INFO    -  Documentation built in 0.14 seconds
```

But by opting for the GitHub link you have the static checking offered by `check-relative-markdown-links`, `:(`.

## Solution (workaround)

Add the relative link within a comment. This way you still get a tripwire for refactors and MkDocs remains happy.

```md
<!--[../testdata/README.md](../testdata/README.md) https://anttiharju.dev/check-relative-markdown-links/comment-trick-explained-->
```

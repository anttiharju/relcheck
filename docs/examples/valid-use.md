# Valid use

This document demonstrates valid use of relative links within markdown as recognized by the `relcheck` tool.

## Links

1. Simple relative links are recognised [Valid use](./valid-use.md)
2. and so are links that traverse upwards [Introduction](../README.md)
3. even files with spaces in their name are supported! See [Issues caught](./issues%20caught.markdown)

## Anchors

1. Anchors can be validated [Introduction#why](../README.md#why)
2. Even duplicate anchors are supported! [Introduction#why-1](../README.md#why-1)

## Code blocks

Markdown links within code blocks are ignored so because they would not be clickable in the rendered document anyway:

```md
[nonexistent](./non.md#existent)
```

## Static check all the things

We can even setup static checking for relative links that we want to have as URLs for whatever reason. Simply add a comment like

```md
<!--[README](./README.md) https://anttiharju.dev/relcheck/comment-trick-explained -->
```

## Image links

![relcheck](../relcheck.png "alt text")

alongside the URL to have the tool detect if the file gets moved in the repo. This makes refactoring project structure a lot less error-prone. Read more about this trick at [https://anttiharju.dev/relcheck/comment-trick-explained](../comment-trick-explained.md)

### Also with single quotes alt text

<!-- prettier-ignore -->
![relcheck](../relcheck.png 'alt text')

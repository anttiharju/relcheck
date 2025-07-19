# Valid use

This document demonstrates valid use of relative links within markdown as recognized by the `relcheck` tool.

## Links

1. Simple relative links are recognised [Valid use](./valid-use.md)
2. and so are links that traverse upwards [Introduction](../README.md)
3. even files with spaces in their name are supported! See [Issues caught](./issues%20caught.markdown)

### With line specified

[like this](./valid-use.md#L5)

## Anchors

1. Anchors can be validated [Introduction#why](../README.md#why)
2. Even duplicate anchors are supported! [Introduction#why-1](../README.md#why-1)

## Code blocks

Markdown links within code blocks are ignored so because they would not be clickable in the rendered document anyway:

```md
[nonexistent](./non.md#existent)
```

<!-- prettier-ignore -->
```also doesn't get confused by exotic formatting like this one```

## ::nut_and_bolt:: Emojis

[Like the one above](./valid-use.md#nut_and_bolt-emojis)

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

<!-- prettier-ignore-start -->
Alternative headings
---
<!-- prettier-ignore-end -->

[alternative headings](./valid-use.md#alternative-headings)

<!-- prettier-ignore-start -->
Alternative headings with equal sign
===
<!-- prettier-ignore-end -->

[alternative headings](./valid-use.md#alternative-headings-with-equal-sign)

# L-starting headings

[l-starting headings](./valid-use.md#L-starting-headings)

<!--[README](./README.md) https://anttiharju.dev/relcheck/comment-trick-explained -->

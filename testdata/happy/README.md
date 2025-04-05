# Happy path tests

This document demonstrates valid use of relative links within markdown as recognized by check-relative-markdown-links tool.

It recognises simple links: [README](./README.md)  
It can follow them upwards: [README](../../README.md)  
It works with files that have spaces in their names: [s p a c e s](./s%20p%20a%20c%20e%20s.md)  
It can handle trivial anchor links: [anchors#hello](./anchors.md#hello)  
It can handle non-trivial anchor links: [anchors#i-have-anchors](./anchors.md#i-have-anchors)  
It can handle duplicate anchor links: [anchors#hello-1](./anchors.md#hello-1)  
It can handle duplicate anchor links: [anchors#hello-2](./anchors.md#hello-2)

It does not care about links within code blocks:

```
[nonexistent](./non.md#existent)
```

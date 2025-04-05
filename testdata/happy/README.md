# Happy path tests

This document demonstrates valid use of relative links within markdown as recognized by check-relative-markdown-links tool.

It recognises simple links: [README](./README.md)  
It can follow them upwards: [README](../../README.md)  
It works with files that have spaces in their names: [s p a c e s](./s%20p%20a%20c%20e%20s.md) <!--%20 is the most compatible way of doing spaces in links afaik-->
It can handle trivial anchor links: [anchors#hello](./anchors.md#hello)  
It can handle non-trivial anchor links: [anchors#i-have-anchors](./anchors.md#i-have-anchors)  
It can handle duplicate anchor links: [anchors#hello-1](./anchors.md#hello-1)  
It can handle duplicate anchor links: [anchors#hello-2](./anchors.md#hello-2)

It does not care about links within code blocks:

```
[nonexistent](./non.md#existent)
```

But it does care about links within comments:

<!--[README](./README.md)-->

This is a trick preserved for MkDocs. MkDocs strictly disallows links outside of the `docs/` directory, in which case you may want to directly link to your documentation site, such as the github version etc. But by bundling in a comment like the one above, you can get static checking for that link via this hidden link. The trick has to be known so it's best to include an explanation with it.

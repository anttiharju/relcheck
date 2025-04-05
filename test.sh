#!/bin/sh
set -e
cd "$(git rev-parse --show-toplevel)/testdata" || exit

../check-relative-markdown-links.bash run --verbose > .got

if [ "$1" = "--regenerate" ]; then
    cp .got want
else
    diff -u want .got
fi

cp want .want

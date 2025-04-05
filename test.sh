#!/bin/sh
cd "$(git rev-parse --show-toplevel)/testdata" || exit

(
    cd happy || exit
    ../../check-relative-markdown-links.bash run --verbose > .got

    if [ "$1" = "--regenerate" ]; then
        cp .got want
    else
        diff --color -u want .got
    fi

    cp want .want
)

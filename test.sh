#!/bin/sh
testdata="$(git rev-parse --show-toplevel)/testdata"
cd "$testdata" || exit

tmp_exit_code="$testdata/.exit_code"
echo 0 > "$tmp_exit_code"

(
    cd happy || exit
    ../../check-relative-markdown-links.bash run --verbose > .got

    if [ "$1" = "--regenerate" ]; then
        cp .got want
    else
        if ! diff --color -u want .got; then
            echo 1 > "$tmp_exit_code"
        fi
    fi

    cp want .want
)

(
    cd negative || exit
    # This style of usage is simple, but not compliant with shellcheck
    # That's the reason why the more complicated usage is bundled in with 'run'
    # shellcheck disable=SC2046
    ../../check-relative-markdown-links.bash --verbose $(git ls-files '*.markdown') > .got

    if [ "$1" = "--regenerate" ]; then
        cp .got want
    else
        if ! diff --color -u want .got; then
            echo 1 > "$tmp_exit_code"
        fi
    fi

    cp want .want
)

exit_code=$(cat "$tmp_exit_code")

exit "$exit_code"

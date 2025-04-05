#!/bin/sh
root="$(git rev-parse --show-toplevel)"
cd "$root" || exit

tmp_exit_code="$root/testdata/.exit_code"
echo 0 > "$tmp_exit_code"

(
    cd docs/examples || exit
    ../../check-relative-markdown-links.bash run --verbose > ../../testdata/got/valid-use
    cd ../../testdata || exit

    if [ "$1" = "--regenerate" ]; then
        cp got/valid-use want/valid-use
    else
        if ! diff --color -u want/valid-use got/valid-use; then
            echo 1 > "$tmp_exit_code"
        fi
    fi
)

(
    cd docs/examples || exit
    ../../check-relative-markdown-links.bash --verbose "$(git ls-files '*.markdown')" > "../../testdata/got/issues caught"
    cd ../../testdata || exit

    if [ "$1" = "--regenerate" ]; then
        cp "got/issues caught" "want/issues caught"
    else
        if ! diff --color -u "want/issues caught" "got/issues caught"; then
            echo 1 > "$tmp_exit_code"
        fi
    fi
)

exit_code=$(cat "$tmp_exit_code")

exit "$exit_code"

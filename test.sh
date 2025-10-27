#!/bin/sh
root="$(git rev-parse --show-toplevel)"
cd "$root" || exit
go build
mkdir -p "$root/tests/got"

tmp_exit_code="$root/tests/.exit_code"
echo 0 > "$tmp_exit_code"

(
    cd docs/examples || exit
    ../../relcheck all --verbose --color=always > ../../tests/got/valid-use
    cd ../../tests || exit

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
    files=$(git ls-files '*.markdown')
    ../../relcheck --verbose --color=always "$files" > "../../tests/got/issues caught"
    cd ../../tests || exit

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

#!/usr/bin/env bash

# https://github.com/anttiharju/check-relative-markdown-links

# Check for relative Markdown links and verify they exist
# Usage: ./check-relative-links.sh [--verbose] file1.md [file2.md] ...
#   or   ./check-relative-links.sh [--verbose] run

set -eu

# Terminal colors and formatting
BOLD="\033[1m"
RED="\033[31m"
YELLOW="\033[33m"
GREEN="\033[32m"
GRAY="\033[90m"
RESET="\033[0m"

# Process arguments
VERBOSE=0
FILES=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose)
            VERBOSE=1
            shift
            ;;
        run)
            # If 'run' argument is provided, use git ls-files to find all markdown files
            # Read each line into the array to be shellcheck compliant
            while IFS= read -r line; do
                FILES+=("$line")
            done < <(git ls-files '*.md')
            shift
            ;;
        *)
            FILES+=("$1")
            shift
            ;;
    esac
done

if [[ ${#FILES[@]} -eq 0 ]]; then
    echo "Usage: $0 [--verbose] <file1.md> [file2.md] ..."
    echo "   or: $0 [--verbose] run  (to check all *.md files in Git)"
    exit 1
fi

EXIT_CODE=0

for file in "${FILES[@]}"; do
    if [[ ! -f "$file" ]]; then
        echo -e "${BOLD}Error:${RESET} ${RED}File not found: ${RESET}$file"
        EXIT_CODE=1
        continue
    fi

    # Get the directory of the current file to resolve relative paths
    dir=$(dirname "$file")

    # Extract all relative links in one pass with awk
    # This avoids multiple grep/awk/echo calls per link
    # Now also capturing the column position
    link_data=$(awk '
        match($0, /\]\(\.[^)]*\)/) {
            link = substr($0, RSTART+2, RLENGTH-3)
            col = RSTART+2  # Column position of the link
            gsub(/#.*$/, "", link)  # Remove anchor part
            if (link != "") {
                print NR ":" col ":" link
            }
        }' "$file")

    # If no links are found, continue to the next file
    if [[ -z "$link_data" ]]; then
        if [[ $VERBOSE -eq 1 ]]; then
            echo -e "${GREEN}✓${RESET} $file: ${GRAY}no relative links${RESET}"
        fi
        continue
    fi

    # Initialize before the subshell
    broken_links_found=0
    valid_links_count=0

    # Process each link
    while IFS=: read -r line_num col_num link; do
        # Construct the full path relative to the file's location
        full_path="$dir/$link"

        if [[ ! -e "$full_path" ]]; then
            # Print the file location in bold
            echo -e "${BOLD}${file}:${line_num}:${col_num}:${RESET} ${RED}broken relative link (file not found):${RESET}"
            # Extract the line content for context
            line_content=$(sed -n "${line_num}p" "$file")
            echo "$line_content"
            # Print line content with yellow indicator pointing to the link position
            printf "${YELLOW}%${col_num}s${RESET}\n" "^"
            broken_links_found=1
        else
            ((valid_links_count++))
        fi
    done < <(echo "$link_data")

    # If verbose mode and we have valid links, report them
    if [[ $VERBOSE -eq 1 && $valid_links_count -gt 0 ]]; then
        if [[ $broken_links_found -eq 0 ]]; then
            if [[ $valid_links_count -eq 1 ]]; then
                echo -e "${GREEN}✓${RESET} $file: found 1 valid relative link"
            else
                echo -e "${GREEN}✓${RESET} $file: found $valid_links_count valid relative links"
            fi
        else
            if [[ $valid_links_count -eq 1 ]]; then
                echo -e "${GRAY}$file: also found 1 valid relative link${RESET}"
            else
                echo -e "${GRAY}$file: also found $valid_links_count valid relative links${RESET}"
            fi
        fi
    fi

    if [[ $broken_links_found -eq 1 ]]; then
        EXIT_CODE=1
    fi
done

# Show a success message if all links are valid, but only in verbose mode
if [[ "$EXIT_CODE" -eq 0 && $VERBOSE -eq 1 ]]; then
    echo -e "${GREEN}✓${RESET} ${BOLD}All relative links are valid!${RESET}"
fi

exit "$EXIT_CODE"

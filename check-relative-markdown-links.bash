#!/usr/bin/env bash
set -u

# https://github.com/anttiharju/check-relative-markdown-links

# Process arguments
verbose=0
force_color=0
files=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose)
            verbose=1
            shift
            ;;
        --color=always)
            force_color=1
            shift
            ;;
        run)
            # If 'run' argument is provided, use git ls-files to find all markdown files
            # Read each line into the array to be shellcheck compliant
            while IFS= read -r line; do
                files+=("$line")
            done < <(git ls-files '*.md')
            shift
            ;;
        *)
            files+=("$1")
            shift
            ;;
    esac
done

if [[ ${#files[@]} -eq 0 ]]; then
    echo "Usage: check-relative-markdown-links [--verbose] [--color=always] <file1.md> [file2.md] ..."
    echo "   or: check-relative-markdown-links [--verbose] [--color=always] run  (to check all *.md files in Git)"
    exit 1
fi

# Terminal colors and formatting
# Default: use colors only when stdout is a terminal
if [ -t 1 ] || [ "$force_color" -eq 1 ]; then
    bold="\033[1m"
    red="\033[31m"
    yellow="\033[33m"
    green="\033[32m"
    gray="\033[90m"
    reset="\033[0m"
else
    # If being piped, use empty strings for colors
    bold=""
    red=""
    yellow=""
    green=""
    gray=""
    reset=""
fi

# Function to URL-decode a string
urldecode() {
    local url_encoded="${1//+/ }"
    printf '%b' "${url_encoded//%/\\x}"
}

# Function to extract all headers from a markdown file and convert to anchors
get_markdown_anchors() {
    local file=$1
    awk '
        # Skip code blocks
        /^```/ {
            in_code_block = !in_code_block
            next
        }
        in_code_block { next }

        # Match headings
        /^#{1,6} / {
            # Extract the heading text without the leading #s
            heading = $0
            gsub(/^#+[ \t]+/, "", heading)
            # Remove trailing spaces
            gsub(/[ \t]+$/, "", heading)

            # Convert to GitHub-style anchor:
            heading = tolower(heading)
            gsub(/[^a-z0-9 -]/, "", heading)
            gsub(/[ \t]+/, "-", heading)
            gsub(/--+/, "-", heading)
            gsub(/-+$/, "", heading)

            # Handle duplicate anchors
            if (anchor_count[heading] > 0) {
                print heading "-" anchor_count[heading]
            } else {
                print heading
            }

            # Increment the counter for this anchor
            anchor_count[heading]++
        }

        BEGIN {
            in_code_block = 0
        }
    ' "$file"
}

exit_code=0

for file in "${files[@]}"; do
    if [[ ! -f "$file" ]]; then
        echo -e "${bold}Error:${reset} ${red}File not found: ${reset}$file"
        exit_code=1
        continue
    fi

    # Get the directory of the current file to resolve relative paths
    dir=$(dirname "$file")

    # Extract all relative links in one pass with awk
    link_data=$(awk '
        # Track if we are inside a code block
        /^```/ {
            in_code_block = !in_code_block
        }

        # Only process links when not in a code block
        !in_code_block && match($0, /\]\(\.[^)]*\)/) {
            link = substr($0, RSTART+2, RLENGTH-3)
            col = RSTART+2  # Column position of the link

            # Keep the original link with anchor for display
            print NR ":" col ":" link
        }

        # Handle triple backtick code blocks without a language specifier
        /^``/ && !match($0, /^```[a-z]*/) {
            in_code_block = !in_code_block
        }

        BEGIN {
            in_code_block = 0
        }
    ' "$file")

    # If no links are found, continue to the next file
    if [[ -z "$link_data" ]]; then
        if [[ $verbose -eq 1 ]]; then
            echo -e "${green}✓${reset} $file: ${gray}no relative links${reset}"
        fi
        continue
    fi

    # Initialize before the subshell
    broken_links_found=0
    valid_links_count=0

    # Process each link
    while IFS=: read -r line_num col_num original_link; do
        # Split the link into path and anchor parts
        if [[ "$original_link" == *"#"* ]]; then
            link_path="${original_link%%#*}"
            link_anchor="${original_link#*#}"
        else
            link_path="$original_link"
            link_anchor=""
        fi

        # URL-decode the link path to handle spaces and other encoded characters
        decoded_link=$(urldecode "$link_path")

        # Construct the full path relative to the file's location
        full_path="$dir/$decoded_link"

        if [[ ! -e "$full_path" ]]; then
            # Print the file location in bold
            echo -e "${bold}${file}:${line_num}:${col_num}:${reset} ${red}broken relative link (file not found):${reset}"
            # Extract the line content for context
            line_content=$(sed -n "${line_num}p" "$file")
            echo "$line_content"
            # Print line content with yellow indicator pointing to the link position
            printf "${yellow}%${col_num}s${reset}\n" "^"
            broken_links_found=1
        elif [[ -n "$link_anchor" ]]; then
            # If an anchor exists, check if it's valid
            target_anchors=$(get_markdown_anchors "$full_path")
            if ! echo "$target_anchors" | grep -Fxq "$link_anchor"; then
                echo -e "${bold}${file}:${line_num}:${col_num}:${reset} ${red}broken relative link (anchor not found):${reset}"
                line_content=$(sed -n "${line_num}p" "$file")
                echo "$line_content"
                printf "${yellow}%${col_num}s${reset}\n" "^"
                broken_links_found=1
            else
                ((valid_links_count++))
            fi
        else
            ((valid_links_count++))
        fi
    done < <(echo "$link_data")

    # If verbose mode and we have valid links, report them
    if [[ $verbose -eq 1 && $valid_links_count -gt 0 ]]; then
        if [[ $broken_links_found -eq 0 ]]; then
            if [[ $valid_links_count -eq 1 ]]; then
                echo -e "${green}✓${reset} $file: found 1 valid relative link"
            else
                echo -e "${green}✓${reset} $file: found $valid_links_count valid relative links"
            fi
        else
            if [[ $valid_links_count -eq 1 ]]; then
                echo -e "${gray}$file: also found 1 valid relative link${reset}"
            else
                echo -e "${gray}$file: also found $valid_links_count valid relative links${reset}"
            fi
        fi
    fi

    if [[ $broken_links_found -eq 1 ]]; then
        exit_code=1
    fi
done

# Show a success message if all links are valid, but only in verbose mode
if [[ "$exit_code" -eq 0 && $verbose -eq 1 ]]; then
    echo -e "${green}✓${reset} ${bold}All relative links are valid!${reset}"
fi

exit "$exit_code"

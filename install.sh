#!/bin/sh
set -eu

# Use mktemp to create a temporary directory
tmp=$(mktemp -d)
# Ensure the temp dir gets cleaned up on exit
trap 'rm -rf "$tmp"' EXIT

# Download the script to the temporary directory
url="https://raw.githubusercontent.com/anttiharju/check-relative-markdown-links/HEAD/check-relative-markdown-links.bash"
curl -o "$tmp/check-relative-markdown-links.bash" "$url"

# Make the script executable
chmod +x "$tmp/check-relative-markdown-links.bash"

# Store the command parts separately
src="$tmp/check-relative-markdown-links.bash"
dest="/usr/local/bin/check-relative-markdown-links"

# Check if sudo needs a password (exit status 1 means password needed)
if ! sudo -n true 2>/dev/null; then
    # Sudo needs password, print the command first
    echo ""
    echo "sudo cp -f $src $dest"
fi

# Run the command with sudo
sudo cp -f "$src" "$dest"

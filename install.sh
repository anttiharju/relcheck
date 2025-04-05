#!/bin/sh
set -eu

# Create the target directory if it doesn't exist
mkdir -p ~/.anttiharju

# Download the script to the target directory
curl -o ~/.anttiharju/check-relative-markdown-links.bash \
  https://raw.githubusercontent.com/anttiharju/check-relative-markdown-links/6d8b08d943582439a074dc0597c081d48fe09243/check-relative-markdown-links.bash

# Make the script executable
chmod +x ~/.anttiharju/check-relative-markdown-links.bash

echo ""
echo "Complete installation by running the following command:"
echo "sudo ln -sf ~/.anttiharju/check-relative-markdown-links.bash /usr/local/bin/check-relative-markdown-links"

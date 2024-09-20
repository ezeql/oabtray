#!/bin/bash

# Check if a version number was provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1
REPO="ezeql/oabtray"
TEMPLATE_FILE="oabtray.rb.template"
TAP_DIR="/opt/homebrew/Library/Taps/ezeql/homebrew-personal"
OUTPUT_FILE="$TAP_DIR/Formula/oabtray.rb"

# Ensure the Formula directory exists
mkdir -p "$TAP_DIR/Formula"

# Download the tar.gz file
TAR_URL="https://github.com/$REPO/archive/refs/tags/v$VERSION.tar.gz"
curl -sL "$TAR_URL" -o "v$VERSION.tar.gz"

# Calculate SHA256
SHA256=$(shasum -a 256 "v$VERSION.tar.gz" | awk '{print $1}')

# Remove the downloaded tar.gz
rm "v$VERSION.tar.gz"

# Replace placeholders in the template and output to the Formula directory
sed -e "s/{{VERSION}}/$VERSION/g" \
    -e "s/{{SHA256}}/$SHA256/g" \
    "$TEMPLATE_FILE" > "$OUTPUT_FILE"

echo "Homebrew formula generated: $OUTPUT_FILE"
echo "Version: $VERSION"
echo "SHA256: $SHA256"
echo "Formula has been placed in your personal tap directory."
echo "To install, run: brew install ezeql/personal/oabtray"
echo "Don't forget to commit and push changes to your tap repository:"
echo "cd $TAP_DIR && git add Formula/oabtray.rb && git commit -m 'Update oabtray to v$VERSION' && git push"
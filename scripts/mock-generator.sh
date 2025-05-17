#!/usr/bin/env bash

set -e

# This script generates mock files for the given Go package.
# It uses the `mockgen` tool to generate the mocks.
#
# Usage:
#  ./scripts/mock-generator.sh <source_dir> <output_dir>

SOURCE_DIR=${1:?"Source directory is required"}
OUTPUT_DIR=${2-"internal/mocks"}

for file in $(find "$SOURCE_DIR" -type f -name "*.go"); do
    # Get the package name from the file path
    package_name=$(basename "$(dirname "$file")")
    # Get the file name without the extension
    file_name=$(basename "$file" .go)
    # Generate the mock file name
    mock_file_name="${file_name}_mock.go"
    # Generate the mock file using mockgen
    echo "Generating mock for $file in $OUTPUT_DIR/$package_name/$mock_file_name for package $package_name"
    mockgen -source="$file" -destination="$OUTPUT_DIR/$package_name/$mock_file_name" -package="${package_name}_mock"
done

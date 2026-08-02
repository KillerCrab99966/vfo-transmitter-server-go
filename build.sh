#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

APP_NAME="vfo-transmitter-server"
OUTPUT_DIR="./bin"
FLAGS="-X 'main.Environment=production'"

# Clean up previous builds
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# Define OS and Architecture matrix
# Syntax: "GOOS/GOARCH"
TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/arm64"
)

echo "Starting cross-compilation..."

for TARGET in "${TARGETS[@]}"; do
  # Split the target string into GOOS and GOARCH
  GOOS=$(echo $TARGET | cut -d'/' -f1)
  GOARCH=$(echo $TARGET | cut -d'/' -f2)

  # Add .exe extension for Windows
  EXTENSION=""
  if [ "$GOOS" = "windows" ]; then
    EXTENSION=".exe"
  fi

  OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}${EXTENSION}"

  echo "Building for $GOOS/$GOARCH -> $OUTPUT_DIR/$OUTPUT_NAME"

  # CGO_ENABLED=0 ensures static linking without native C dependencies
  # -ldflags="-s -w" strips debug information to make file sizes smaller
  CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
    -trimpath -ldflags="-s -w $FLAGS" \
    -o "$OUTPUT_DIR/$OUTPUT_NAME" .
done

echo "Done! Artifacts are in $OUTPUT_DIR"

#!/bin/bash

VERSION=$(git rev-parse --short HEAD)

echo "Version: $VERSION"

go build \
    -ldflags="-X 'main.VERSION=$VERSION'" \
    -o oabtray

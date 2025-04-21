#!/usr/bin/env bash
set -euo pipefail

# You can replace the image path with a path to your own image file.
# You must replace the JWT with a valid JWT token from the output of the server.sh script
# You can replace the tags with your own tags, separated by commas.

go run ./cmd/client \
-img-path "demo-image.jpeg" \
-tags "" \
-jwt ""

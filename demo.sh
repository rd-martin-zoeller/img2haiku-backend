#!/usr/bin/env bash
set -euo pipefail

set -a
source .env.debug
set +a

go run ./cmd

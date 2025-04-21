#!/usr/bin/env bash
set -euo pipefail

export LOCAL_ONLY=true
export FUNCTION_TARGET=ComposeHaiku
export OPENAI_API_KEY=sk-... #Replace with your OpenAI API key

go run ./cmd/server

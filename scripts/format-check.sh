#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
go_diff="$(find "$repo_root/core" "$repo_root/webapp" -name '*.go' -print0 | xargs -0 gofmt -d)"

if [[ -n "$go_diff" ]]; then
    printf '%s\n' "$go_diff"
    exit 1
fi

(
    cd "$repo_root/frontend"
    npm run format:check
)
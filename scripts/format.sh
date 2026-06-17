#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

find "$repo_root/core" "$repo_root/webapp" -name '*.go' -print0 | xargs -0 gofmt -w

(
    cd "$repo_root/frontend"
    npm run format
)
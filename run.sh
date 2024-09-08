#!/bin/sh

set -e # Exit early if any commands fail


(
  cd "$(dirname "$0")"
  go build -o /tmp/interpreter-target ./cmd/myinterpreter
)
exec /tmp/interpreter-target "$@"

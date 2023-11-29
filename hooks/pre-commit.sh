#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
if [[ "${TRACE-0}" == "1" ]]; then
  set -o xtrace
fi

cd "$(dirname "$0")"

main() {
  echo "Running pre-commit hook"

  (
    cd "../../" &&
	echo "Running all tests" && go test ./...
  )

  echo "Finished pre-commit hook"
}

main "$@"

#! /bin/bash

set -e

FILES="$(git diff --cached --name-only --diff-filter=ACMR | sed 's| |\\ |g' | grep -i '.*\.go$' || true)"

[[ -z "$FILES" ]] && exit 0

# Organize imports with goimports
echo "Running goimports"
echo "$FILES" | xargs -n 1 goimports -w

# Prettify all selected files
echo "Running golines"
echo "$FILES" | xargs -n 1 golines --base-formatter="gofmt" --write-output

# Add back the modified/prettified files to staging
echo "$FILES" | xargs git add

echo "Running golangci-lint"
golangci-lint run -v

exit 0

#! /bin/bash

set -e

MINIMUM_VERSION_FOR_GOLINES_LASTED='1.19'

go_version=$(go version | awk '{print $3}' | sed 's/go//')
golines_download_cmd='go install github.com/segmentio/golines@v0.9.0'

log_info() {
  echo -e "\033[0;32m[INFO]\033[0m: $1"
}

log_warn() {
  echo -e "\033[1;33m[WARN]\033[0m: $1"
}

log_error() {
  echo -e "\033[0;31m[ERRO]\033[0m: $1"
}

install_goimports() {
  if [[ -x "$(command -v goimports)"  ]]; then
    log_info 'goimports have already been installed'
    return
  fi

  echo 'goimports not found, install using go'
  bash -c 'go install golang.org/x/tools/cmd/goimports@latest'
}

install_golines() {
  if [[ -x "$(command -v golines)"  ]]; then
    log_info 'golines have already been installed'
    return
  fi

  log_info 'golines not found, install using go'

  if [[ "$(printf '%s\n' "$go_version" "$MINIMUM_VERSION_FOR_GOLINES_LASTED" | sort -V | head -n1)" = "$MINIMUM_VERSION_FOR_GOLINES_LASTED" ]]; then
    golines_download_cmd="go install github.com/segmentio/golines@latest"
  fi

  log_info "golines download command: $golines_download_cmd"

  log_info "Install golines"
  bash -c "$golines_download_cmd"

  log_info "Checkout https://github.com/segmentio/golines#developer-tooling-integration for IDE integration"
}

install_golangci_lint() {
  if [[ -x "$(command -v golangci-lint)"  ]]; then
    log_info 'golangci-lint have already been installed'
    return
  fi

  log_info 'golangci-lint not found'

  if [[ "$(uname -s)" == "Darwin" ]]; then
    log_info 'Machine is Darwin'
    if [[ $commands[brew] ]]; then
      log_info 'brew present, install golangci_lint using brew'
      bash -c 'brew install golangci-lint'
      log_info "Checkout https://golangci-lint.run/usage/integrations/ for IDE integration"
    else
      log_error 'brew not present, checkout https://golangci-lint.run/usage/install for alternative installation'
    fi
  else
    log_error 'Machine is not Darwin, checkout https://golangci-lint.run/usage/install for alternative installation'
  fi
}

create_pre_commit_hook() {
  if [[ -f $(pwd)/.git/hooks/pre-commit ]]; then
    log_warn 'Pre-commit hook file is already exist, file will not be copy'
    log_warn 'Delete the old pre-commit file and run this script again or add the content of ./scripts/pre-commit.sh manually'
  else
    log_info "Copy file ./scripts/pre-commit.sh to .git/hooks/pre-commit"
    cp ./scripts/pre-commit.sh ./.git/hooks/pre-commit
    log_info "Grant execute permission to the pre-commit file"
    chmod +x ./.git/hooks/pre-commit
  fi
}

main() {
  log_info 'Init formatting and linting tools'
  log_info "Go version $go_version"

  install_goimports
  install_golines
  install_golangci_lint

  create_pre_commit_hook
}


main

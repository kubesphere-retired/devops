#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

echo $(go version)

ROOT=$(dirname "${BASH_SOURCE}")/..
source "${ROOT}/hack/util.sh"
cd "${ROOT}"

type goimports >/dev/null 2>&1|| { echo "goimports is required. Please install it by `go get golang.org/x/tools/cmd/goimports`"; exit 1;}

if [ $# -gt 0 ]; then
  case $1 in
    -d|--diff)
      echo "Verify files changed......"
      bad_files=$( util::find_diff_files| xargs goimports -l  -e -local=kubesphere )
      ;;
    -a|--all)
      echo "Verify all files......"
      bad_files=$(util::find_files| xargs goimports -l  -e -local=kubesphere )
      ;;
    *)
      echo "Please specify -a (verify all go files) or -d (verify diff only)"
      exit 1
      ;;
  esac
else
  echo "Please specify -a (verify all go files) or -d (verify diff only)"
  exit 1
fi
bad_files=$(util::find_files| xargs goimports -l  -e -local=kubesphere )
if [[ -n "${bad_files}" ]]; then
  echo "!!! The imports of following files should be formated: " >&2
  echo "${bad_files}"
  echo "Try running 'goimports -l -w -e -local=kubesphere [path]'" >&2
  exit 1
fi
echo "goimports verify PASS"

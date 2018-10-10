#!/bin/bash

set -o nounset
set -o pipefail
set -o errexit
echo $(go version)

ROOT=$(dirname "${BASH_SOURCE}")/..
source "${ROOT}/hack/util.sh"
cd "${ROOT}"

type goimports >/dev/null 2>&1|| { echo "goimports is required. Please install it by `go get golang.org/x/tools/cmd/goimports`"; exit 1;}

set +e
if [ $# -gt 0 ]; then
  case $1 in
    -d|--diff)
      echo "Verify files changed......"
   
      bad_files=$(util::find_diff_files)
      set -e
      ;;
    -a|--all)
      echo "Verify all files......"
      bad_files=$(util::find_files)
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
set -e

if [ -z "$bad_files" ]; then
    echo "No go files need to check. PASS"
    exit 0
fi

bad_files=$(echo -e $bad_files| xargs goimports -l  -e -local=kubesphere )
if [[ -n "${bad_files}" ]]; then
  echo "The imports of following files should be formated: " 
  echo "${bad_files}"
  echo "Try running 'goimports -l -w -e -local=kubesphere [path]'" 
  exit 1
fi
echo "goimports verify PASS"


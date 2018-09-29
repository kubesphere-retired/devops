#!/bin/bash

set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE}")/..
source "${ROOT}/hack/util.sh"
cd "${ROOT}"

type goimports >/dev/null 2>&1|| { echo "goimports is required. Please install it by `go get golang.org/x/tools/cmd/goimports`"; exit 1;}

USAGE="Please specify -a (format all go files) or -d (format files changed between current commit and last commit) or -l (format files changed locally)"
if [ $# -gt 0 ]; then
  case $1 in
    -d|--diff)
      echo "Format files between commits......."
      bad_files=$(util::find_diff_files)
      ;;
    -a|--all)
      echo "Format all files......"
      bad_files=$(util::find_files)
      ;;
    -l|--local)
      echo "Format local changed files......"
      bad_files=$(util::find_local_change_files)
      ;;
    *)
      echo "$USAGE"
      exit 1
      ;;
  esac
else
  echo "$USAGE"
  exit 1
fi

if [ -z $bad_files ]; then
    echo "No go files need to format. DONE"
    exit 0
fi
for file in $bad_files
do
    echo "Processing ::$file"
    goimports -l  -w -e -local=kubesphere  $file
done

if [ $? -ne 0 ]; then
    echo "Cannot to format source codes, Please check it manualy "
    exit 1
fi
echo "Format is DONE"


#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE}")
cd $ROOT

./verify-goimports.sh -a

echo "ALL PASS"


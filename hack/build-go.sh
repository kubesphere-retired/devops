#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname $BASH_SOURCE)/..

cd $ROOT

CGO_ENABLED=0 GOOS=linux go build -v  -a -installsuffix cgo -ldflags '-w'  -o cmd/server pkg/server/main.go


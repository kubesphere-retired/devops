#!/bin/bash


set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname $BASH_SOURCE)/..

cd $ROOT
if [ $# -lt 1 ]; then
    echo "Please enter the name of image"
    exit 1
fi  

IMAGE_NAME=$1
docker build -t $IMAGE_NAME -f build/docker/Dockerfile .



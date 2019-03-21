#!/bin/bash

Version=$1

make build-dev-image-${Version}
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
make push-dev-image-${Version}

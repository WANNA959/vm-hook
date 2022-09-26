#!/bin/bash

set -x

export DOCKER_USER=wannazjx
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn,direct

CURRENT_DIR=$(cd "$(dirname "$0")";pwd)
# build webhook
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o label-hook $CURRENT_DIR/../.
# build docker image
docker build --no-cache -f $CURRENT_DIR/Dockerfile -t ${DOCKER_USER}/label-hook:v1 .
rm -rf label-hook

docker push ${DOCKER_USER}/label-hook:v1

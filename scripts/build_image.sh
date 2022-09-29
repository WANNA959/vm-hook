#!/bin/bash

set -x

export DOCKER_HUB_URL=registry.cn-beijing.aliyuncs.com
export IMAGE_TAG_PREFIX=${DOCKER_HUB_URL}/dosstack
export DOCKER_USER=netgenius201
#export DOCKER_USER=wannazjx
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn,direct

CURRENT_DIR=$(cd "$(dirname "$0")";pwd)
# build webhook
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o label-hook $CURRENT_DIR/../.
# build docker image
docker buildx create --name mybuilder --driver docker-container
docker buildx use mybuilder
docker run --privileged --rm tonistiigi/binfmt --install all

docker login -u ${DOCKER_USER} ${DOCKER_HUB_URL}
#docker login -u ${DOCKER_USER}
docker buildx build --no-cache -f $CURRENT_DIR/Dockerfile -t ${IMAGE_TAG_PREFIX}/kubestack-label-hook:v1 --platform linux/arm64,linux/amd64 . --push
rm -rf label-hook


#docker push ${DOCKER_USER}/label-hook:v1

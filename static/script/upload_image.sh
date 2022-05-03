#!/bin/zsh

USER_NAME="g3un"
IMAGE_NAME="sp1-backend"
IMAGE_TAG="v0.0.6"

docker build -t ${USER_NAME}/${IMAGE_NAME}:${IMAGE_TAG} . &&\
docker push ${USER_NAME}/${IMAGE_NAME}:${IMAGE_TAG} &&\
docker image rm ${USER_NAME}/${IMAGE_NAME}:${IMAGE_TAG}

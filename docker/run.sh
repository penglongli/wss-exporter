#!/bin/bash

IMPORT_PATH="wss-exporter"
APP_NAME="wss-exporter"

rm -rf bin/

# go build
docker run --rm \
    -v ${PWD%/*}:/app/src/${IMPORT_PATH} \
    -w /app \
    -e GOPATH=/app \
    golang:1.9.7 \
    go build -o src/${IMPORT_PATH}/docker/bin/${IMPORT_PATH} ${IMPORT_PATH}

# image build
docker build -t ${APP_NAME} .

# image run
docker rm -f ${APP_NAME}
docker run -it -d \
    -p 8080:8080 \
    --name=${APP_NAME} \
    ${APP_NAME}
    

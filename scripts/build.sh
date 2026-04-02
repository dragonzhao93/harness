#!/bin/sh

export GOOS=linux
export GOARCH=amd64

echo "building docker images for ${GOOS}/${GOARCH} ..."

REPO="github.com/drone/drone"

# compile the server using the cgo
go build -ldflags "-extldflags \"-static\"" -o release/linux/${GOARCH}/drone-server ${REPO}/cmd/drone-server

docker build --no-cache -f ./docker/Dockerfile.server.${GOOS}.${GOARCH} -t drone/drone:2.28.1-huifu . --platform ${GOOS}/${GOARCH}

#!/bin/bash

docker pull quay.io/goswagger/swagger

for f in $(find ./cmd -name *.go)
do
    docker run --rm -it -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(realpath $(dirname $f)) quay.io/goswagger/swagger generate spec --scan-models -o swagger.json
done
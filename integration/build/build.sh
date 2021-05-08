#!/bin/bash

# Usage: 
#    integration/build/build.sh
#    integration/build/build.sh force # Always recreate docker image and container.
set -e

SCRIPTPATH=$(dirname "$0")

echo $SCRIPTPATH

if [ "$1" =  "force" ] || [[ "$(docker images -q dnero_builder_image 2> /dev/null)" == "" ]]; then
    docker build -t dnero_builder_image $SCRIPTPATH
fi

set +e
docker stop dnero_builder
docker rm dnero_builder
set -e

docker run --name dnero_builder -it -v "$GOPATH:/go" dnero_builder_image

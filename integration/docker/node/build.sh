#!/bin/bash

# Build a docker image for a Dnero node.
# Usage: 
#    integration/docker/node/build.sh
#
# After the image is built, you can create a container by:
#    docker stop dnero_node
#    docker rm dnero_node
#    docker run -e DNERO_CONFIG_PATH=/dnero/integration/privatenet/node --name dnero_node -it dnero
set -e

SCRIPTPATH=$(dirname "$0")

echo $SCRIPTPATH

if [ "$1" =  "force" ] || [[ "$(docker images -q dnero 2> /dev/null)" == "" ]]; then
    docker build -t dnero -f $SCRIPTPATH/Dockerfile .
fi



#!/bin/bash

echo $1
if [ "$1" = "-r" ]; then
    echo rebuilding
    docker build --no-cache -t docker_default .
fi
docker-compose up -d

url=$(docker-compose exec jupyter jupyter notebook list | grep http | awk '{print $1}')
if [[ -z $url ]]; then
    echo Cannot determine url
    exit 1
fi

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    xdg-open $url
elif [[ "$OSTYPE" == "darwin"* ]]; then
    open $url
else
    echo $url
fi

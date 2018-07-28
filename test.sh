#!/bin/sh 
set -e

docker-compose build
docker-compose run --rm --no-deps api go test -v